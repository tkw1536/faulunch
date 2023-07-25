package faulunch

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/swaggest/swgui/v4emb"

	_ "embed"
)

//go:embed openapi.json
var openAPIJSON []byte

func (server *Server) handleAPI() http.Handler {
	// build the swagger api
	swagger := v4emb.NewHandler("FauLunch API", "/api/openapi.json", "/api/")

	// build the router api v1
	v1router := httprouter.New()
	v1router.Handle(http.MethodGet, "/api/v1/sync", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		server.handleAPISync(w, r)
	})
	v1router.Handle(http.MethodGet, "/api/v1/locations", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		server.handleAPILocations(w, r)
	})
	v1router.Handle(http.MethodGet, "/api/v1/menu/:location", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		location := Location(p.ByName("location"))
		server.handleAPIMenuDays(location, w, r)
	})
	v1router.Handle(http.MethodGet, "/api/v1/menu/:location/:day", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		day := ParseDay(p.ByName("day"))
		location := Location(p.ByName("location"))
		server.handleAPIMenu(location, day, w, r)
	})

	var mux http.ServeMux
	mux.Handle("/api/", swagger)
	mux.Handle("/api/v1/", v1router)
	mux.Handle("/api/openapi.json", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openAPIJSON)
	}))
	return &mux
}

const notFoundError = `{"status":"Not Found"}`

// handleNotFound sends a not found response to the caller
func (server *Server) handleNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(notFoundError))
}

const (
	internalServerError = `{"status":"Internal Server Error"}`
)

// handleInternalServerError sends an internal server error
func (server *Server) handleInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(internalServerError))
}

func (server *Server) handleAPISync(w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.Sync").Logger()

	// fetch all the items
	sync, err := server.API.LastSync()
	logger.Trace().Err(err).Msg("API.Sync")

	if err != nil {
		server.handleInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sync)
}

func (server *Server) handleAPILocations(w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.Locations").Logger()

	// fetch all the items
	results, err := server.API.Locations()
	logger.Trace().Err(err).Msg("API.Locations")

	if err != nil || len(results) == 0 {
		server.handleInternalServerError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (server *Server) handleAPIMenuDays(location Location, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.MenuDays").Str("location", string(location)).Logger()

	// get the day to start with
	from := ParseDay(r.URL.Query().Get("from"))
	if from == 0 {
		from = Today().Add(-21) // default 21 days ago
	}

	count, err := strconv.Atoi(r.URL.Query().Get("days"))
	if err != nil {
		count = 28 // default to 28 days
	}

	// clamp to the range 1 ... 365
	if count < 1 {
		count = 1
	}
	if count > 365 {
		count = 365
	}

	results, err := server.API.Days(location, from, count)
	logger.Trace().Err(err).Msg("API.Days")

	// something went wrong
	if err != nil {
		server.handleInternalServerError(w)
		return
	}

	// check if the location exists
	if len(results) == 0 {
		exists, err := server.API.KnowsLocation(location)
		if err != nil {
			server.handleInternalServerError(w)
			return
		}

		if !exists {
			server.handleNotFound(w)
			return
		}
	}

	// send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (server *Server) handleAPIMenu(location Location, day Day, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.Menu").Str("location", string(location)).Stringer("day", day).Logger()

	results, err := server.API.MenuItems(location, day)
	logger.Trace().Err(err).Msg("API.MenuItems")

	if err != nil {
		server.handleInternalServerError(w)
		return
	}

	if len(results) == 0 {
		server.handleNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
