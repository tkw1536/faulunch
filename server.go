package faulunch

import (
	"embed"
	"html/template"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"github.com/tkw1536/faulunch/internal"
	"golang.org/x/text/language"
)

//go:embed "api_server"
var apiServerData embed.FS
var apiServerTemplate = template.Must(template.ParseFS(apiServerData, "api_server/*.html", "api_server/index.js", "api_server/index.css"))

type Server struct {
	API    *API
	Logger *zerolog.Logger

	init   sync.Once
	router *httprouter.Router
	Legal  ServerLegal
}

type ServerLegal struct {
	Link     string
	DEString string
	ENString string
}

func (legal ServerLegal) DEHTML() template.HTML {
	if legal.Link == "" || legal.DEString == "" {
		return ""
	}
	return template.HTML("<a href='" + legal.Link + "' target='_blank' rel='noopener noreferer'>" + template.HTMLEscapeString(legal.DEString) + "</a>")
}
func (legal ServerLegal) ENHTML() template.HTML {
	if legal.ENString == "" {
		legal.ENString = legal.DEString
	}
	if legal.Link == "" || legal.ENString == "" {
		return ""
	}
	return template.HTML("<a href='" + legal.Link + "' target='_blank' rel='noopener noreferer'>" + template.HTMLEscapeString(legal.ENString) + "</a>")
}

var matcher = language.NewMatcher([]language.Tag{
	language.German,
	language.English,
})

var regexpValidLocation = regexp.MustCompile("^[aA-zZ0-9-_]+$")

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.init.Do(func() {
		server.router = httprouter.New()

		server.router.Handle(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			tag, _ := language.MatchStrings(matcher, r.Header.Get("Accept-Language"))
			if tag == language.German {
				http.Redirect(w, r, "/de/", http.StatusTemporaryRedirect)
			} else {
				http.Redirect(w, r, "/en/", http.StatusTemporaryRedirect)
			}
		})

		// index
		server.router.Handle(http.MethodGet, "/en/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			server.HandleIndex(true, w, r)
		})
		server.router.Handle(http.MethodGet, "/de/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			server.HandleIndex(false, w, r)
		})

		// specific locations
		locations, _ := server.API.Locations()
		for _, l := range locations {
			l := string(l)
			if !regexpValidLocation.MatchString(l) || l == "de" || l == "en" {
				server.Logger.Warn().Str("location", l).Msg("Skipping invalid location")
			}
			server.router.Handle(http.MethodGet, "/"+l+"/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
				tag, _ := language.MatchStrings(matcher, r.Header.Get("Accept-Language"))
				if tag == language.German {
					http.Redirect(w, r, "/de/"+l+"/", http.StatusTemporaryRedirect)
				} else {
					http.Redirect(w, r, "/en/"+l+"/", http.StatusTemporaryRedirect)
				}
			})
		}

		// location
		server.router.Handle(http.MethodGet, "/en/:location/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			location := Location(p.ByName("location"))
			server.HandleLocation(location, true, w, r)
		})

		server.router.Handle(http.MethodGet, "/de/:location/", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			location := Location(p.ByName("location"))
			server.HandleLocation(location, false, w, r)
		})

		// menu

		server.router.Handle(http.MethodGet, "/en/:location/:day", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			day := ParseDay(p.ByName("day"))
			location := Location(p.ByName("location"))
			server.HandleMenu(location, day, true, w, r)
		})

		server.router.Handle(http.MethodGet, "/de/:location/:day", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
			day := ParseDay(p.ByName("day"))
			location := Location(p.ByName("location"))
			server.HandleMenu(location, day, false, w, r)
		})

		// API
		server.router.Handler(http.MethodGet, "/api/*filepath", server.handleAPI())
	})

	server.router.ServeHTTP(w, r)
}

type globalContext struct {
	requestURI string // uri being requested
	English    bool

	legal    ServerLegal
	LastSync time.Time
}

func (gc *globalContext) loadLastSync(api *API) error {
	lastSync, err := api.LastSync()
	if err != nil {
		return err

	}
	gc.LastSync = time.Unix(lastSync.Stop, 0).UTC()
	return nil
}

func (gc globalContext) LegalHTML() template.HTML {
	if gc.English {
		return gc.legal.ENHTML()
	} else {
		return gc.legal.DEHTML()
	}
}

func (gc globalContext) LangAttr() template.HTMLAttr {
	if gc.English {
		return "lang=\"en\""
	}
	return "lang=\"de\""
}

func (gc globalContext) Alternate() template.HTML {
	alternate := gc.requestURI
	var lang, title string
	if gc.English {
		alternate = strings.ReplaceAll(alternate, "/en/", "/de/")
		lang, title = "de", "🇩🇪 Deutsche Version"
	} else {
		alternate = strings.ReplaceAll(alternate, "/de/", "/en/")
		lang, title = "en", "🇬🇧 English Version"
	}

	return template.HTML("<a href='" + alternate + "' rel='alternate' lang='" + lang + "'>" + title + "</a>")
}

type indexContext struct {
	globalContext
	Locations []Location
}

func (server *Server) HandleIndex(english bool, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "Index").Logger()

	// fetch all the items
	results, err := server.API.Locations()
	logger.Debug().Err(err).Msg("API.Locations")

	if err != nil || len(results) == 0 {
		http.NotFound(w, r)
		return
	}

	// and execute the template
	{
		context := indexContext{
			globalContext: globalContext{
				English:    english,
				requestURI: r.URL.RequestURI(),
				legal:      server.Legal,
			},
			Locations: results,
		}
		if err := context.loadLastSync(server.API); err != nil {
			logger.Debug().Err(err).Msg("LoadLastSync")
		}

		w.Header().Add("Content-Type", "text/html")
		err := apiServerTemplate.ExecuteTemplate(w, "index.html", context)
		logger.Debug().Err(err).Msg("ExecuteTemplate")
	}
}

type menuContext struct {
	globalContext

	Day        Day
	Location   Location
	Pagination Pagination
	Items      []MenuItem

	Allergens   []Allergen
	Additives   []Additive
	Ingredients []Ingredient
}

func (mc menuContext) ID(id string) string {
	return strings.ReplaceAll(id, " ", "-")
}

func (mc menuContext) Link(d Day) template.HTML {
	link := string(mc.Location) + "/" + d.String()
	var date string
	if mc.globalContext.English {
		link = "/en/" + link
		date = string(d.ENHTML())
	} else {
		link = "/de/" + link
		date = string(d.DEHTML())
	}

	return template.HTML("<a href='" + link + "'>" + date + "</a>")
}

const (
	menuPaginationSize = 2
)

func (server *Server) HandleLocation(location Location, english bool, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "HandleLocation").Str("location", string(location)).Logger()

	now, err := server.API.CurrentDay(location, Today())
	logger.Debug().Err(err).Msg("API.CurrentDay")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	server.HandleMenu(location, now, english, w, r)
}

func (server *Server) HandleMenu(location Location, day Day, english bool, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "HandleMenu").Str("location", string(location)).Stringer("day", day).Logger()

	mc := menuContext{
		globalContext: globalContext{
			English:    english,
			requestURI: r.URL.RequestURI(),
			legal:      server.Legal,
		},
		Location: location,
		Day:      day,
	}

	if err := mc.loadLastSync(server.API); err != nil {
		logger.Debug().Err(err).Msg("LoadLastSync")
	}

	var err error

	// fetch all the items
	mc.Items, err = server.API.MenuItems(location, day)
	logger.Debug().Err(err).Msg("API.MenuItems")
	if err != nil || len(mc.Items) == 0 {
		http.NotFound(w, r)
		return
	}

	mc.Pagination, err = server.API.DayPagination(location, day, menuPaginationSize)
	logger.Debug().Err(err).Msg("API.DayPagination")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// merge all the annotations
	additivesSet := make(map[Additive]struct{})
	allergensSet := make(map[Allergen]struct{})
	ingredientsSet := make(map[Ingredient]struct{})

	for _, i := range mc.Items {
		for _, add := range i.AdditiveAnnotations.Data {
			additivesSet[add] = struct{}{}
		}
		for _, allergen := range i.AllergenAnnotations.Data {
			allergensSet[allergen] = struct{}{}
		}
		for _, ing := range i.IngredientAnnotations.Data {
			ingredientsSet[ing] = struct{}{}
		}
	}

	mc.Additives = internal.SortedKeysOf(additivesSet, func(a, b Additive) int { return a.Cmp(b) })
	mc.Allergens = internal.SortedKeysOf(allergensSet, func(a, b Allergen) int { return a.Cmp(b) })
	mc.Ingredients = internal.SortedKeysOf(ingredientsSet, func(a, b Ingredient) int { return a.Cmp(b) })

	// and execute the template
	{
		w.Header().Add("Content-Type", "text/html")
		err := apiServerTemplate.ExecuteTemplate(w, "menu.html", mc)
		logger.Debug().Err(err).Msg("ExecuteTemplate")
	}
}
