package router

import (
	"net/http"

	"github.com/LeReverandNox/GuessWhat/src/handlers"
	"github.com/gorilla/mux"
)

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/", handlers.IndexHandler).Methods("GET")
}

// NewRouter is a method that encapsulate a Mux router, with all it's routes
func NewRouter() http.Handler {
	router := mux.NewRouter()

	registerRoutes(router)

	return router
}
