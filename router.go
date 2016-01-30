package main

import (
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

	router.HandleFunc("/login", LogIn).Methods("POST")

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
