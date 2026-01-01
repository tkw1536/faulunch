//spellchecker:words faulunch
package faulunch

//spellchecker:words encoding json http strconv github swaggest swgui embed
import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/swaggest/swgui/v5emb"
	"github.com/tkw1536/faulunch/internal/location"
	"github.com/tkw1536/faulunch/internal/ltime"

	_ "embed"
)

//go:embed openapi.json
var openAPIJSON []byte

// registerAPIRoutes registers API routes to the server mux
func (server *Server) registerAPIRoutes() {
	// api + documentation
	server.mux.Handle("GET /api/", v5emb.NewHandler("FauLunch API", "/api/openapi.json", "/api/"))
	server.mux.HandleFunc("GET /api/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(openAPIJSON)
	})

	// api endpoints
	server.mux.HandleFunc("GET /api/v1/healthcheck", server.handleAPIHealth)
	server.mux.HandleFunc("GET /api/v1/sync", server.handleAPISync)
	server.mux.HandleFunc("GET /api/v1/locations", server.handleAPILocations)
	server.mux.HandleFunc("GET /api/v1/menu/{location}", server.handleAPIMenuDays)
	server.mux.HandleFunc("GET /api/v1/menu/{location}/{day}", server.handleAPIMenu)
	server.mux.HandleFunc("GET /api/v1/sqlite", server.handleAPIsqlite)
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
	statusHealthy       = `{"status":"healthy"}`
)

// handleInternalServerError sends an internal server error
func (server *Server) handleInternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(internalServerError))
}

func (server *Server) handleAPIHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db, err := server.API.DB.DB()
	if err != nil || db.PingContext(r.Context()) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(statusHealthy))
}

func (server *Server) handleAPISync(w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.Sync").Logger()

	// fetch all the items
	sync, err := server.API.LastSync(r.Context())
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

func (server *Server) handleAPIMenuDays(w http.ResponseWriter, r *http.Request) {
	location := location.Location(r.PathValue("location"))

	logger := server.Logger.With().Str("route", "API.MenuDays").Str("location", string(location)).Logger()

	// get the day to start with
	from := ltime.ParseDay(r.URL.Query().Get("from"))
	if from == 0 {
		from = ltime.Today().Add(-21) // default 21 days ago
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

func (server *Server) handleAPIMenu(w http.ResponseWriter, r *http.Request) {
	day := ltime.ParseDay(r.PathValue("day"))
	location := location.Location(r.PathValue("location"))

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

func (server *Server) handleAPIsqlite(w http.ResponseWriter, r *http.Request) {
	if server.API.Copier == nil {
		server.handleNotFound(w)
		return
	}

	logger := server.Logger.With().Str("route", "API.sqlite").Logger()

	err := server.API.Copier(w, r)
	logger.Trace().Err(err).Msg("API.sqlite")

	if err != nil {
		server.handleInternalServerError(w)
		return
	}
}
