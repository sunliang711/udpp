package handlers

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sunliang711/udpp/types"
	"net/http"
)

//Register all handlers
var allHandlers = []*types.HandlerObj{
	{"/login", http.HandlerFunc(loginHandler), []string{"POST"}, ""},
	{"/config", auth(http.HandlerFunc(getConfig)), []string{"POST"}, ""},
	{"/update_config", auth(http.HandlerFunc(updateConfig)), []string{"POST"}, ""},
	{"/get_rights", auth(http.HandlerFunc(getRights)), []string{"POST"}, ""},
	{"/update_rights", auth(http.HandlerFunc(updateRights)), []string{"POST"}, ""},
}

func Router(enableCors bool) http.Handler {
	rt := mux.NewRouter()

	for _, h := range allHandlers {
		rt.Handle(h.Path, h.H).Methods(h.Methods...)
	}

	if enableCors {
		return cors.Default().Handler(rt)
	} else {
		return rt
	}
}
