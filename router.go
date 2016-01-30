package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	tokenAuth := NewTokenAuth(nil, nil, memStore, nil)

	router.HandleFunc("/login/{id}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		t := memStore.NewToken(vars["id"])
		fmt.Fprintf(w, "hi %s, your token is %s", vars["id"], t)
	})

	for _, route := range routes {
		var handler http.Handler

		handler = tokenAuth.HandleFunc(route.HandlerFunc)
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
