package faulunch

import (
	"encoding/json"
	"net/http"

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

const (
	internalServerError = `{"status":"Internal Server Error"}`
	notFoundError       = `{"status":"Not Found"}`
)

func (server *Server) handleAPISync(w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.Sync").Logger()

	// fetch all the items
	sync, err := server.API.LastSync()
	logger.Trace().Err(err).Msg("API.Sync")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (server *Server) handleAPIMenuDays(location Location, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.MenuDays").Str("location", string(location)).Logger()

	results, err := server.API.Days(location)
	logger.Trace().Err(err).Msg("API.Days")

	if err != nil || len(results) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(notFoundError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func (server *Server) handleAPIMenu(location Location, day Day, w http.ResponseWriter, r *http.Request) {
	logger := server.Logger.With().Str("route", "API.Menu").Str("location", string(location)).Stringer("day", day).Logger()

	results, err := server.API.MenuItems(location, day)
	logger.Trace().Err(err).Msg("API.MenuItems")

	if err != nil || len(results) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(notFoundError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
