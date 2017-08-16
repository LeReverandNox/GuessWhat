package routing

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/LeReverandNox/GuessWhat/src/handlers"
	"github.com/gorilla/mux"
)

type MyRouter mux.Router

func registerRoutes(router *mux.Router) {
	router.HandleFunc("/", handlers.IndexHandler).Methods("GET")

	router.PathPrefix("/css").Handler(http.StripPrefix("/css", http.FileServer(http.Dir("assets/css/"))))
	router.PathPrefix("/js").Handler(http.StripPrefix("/js", http.FileServer(http.Dir("assets/js/"))))
}

// NewRouter is a method that encapsulate a Mux router, with all it's routes
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	registerRoutes(router)

	return router
}

func ShowRoutes(router *mux.Router) {
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		// p will contain a regular expression that is compatible with regular expressions in Perl, Python, and other languages.
		// For example, the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'.
		p, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(m, ","), t, p)
		return nil
	})
}
