package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sunliang711/udpp/types"
)

//Router returns router handler
func Router(enableCors bool) http.Handler {
	//Register all handlers
	var allHandlers = []*types.HandlerObj{
		{"/login", http.HandlerFunc(loginHandler), []string{"POST"}, ""},
		{"/get_config", auth(http.HandlerFunc(getConfig)), []string{"GET"}, ""},
		{"/update_config", auth(http.HandlerFunc(updateConfig)), []string{"POST"}, ""},
		{"/get_right", auth(http.HandlerFunc(getRights)), []string{"GET"}, ""},
		{"/update_right", auth(http.HandlerFunc(updateRights)), []string{"POST"}, ""},
		{"/access_list", auth(http.HandlerFunc(accessList)), []string{"GET"}, ""},
	}
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
