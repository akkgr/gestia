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

	tokenAuth := NewTokenAuth(nil, nil, memStore, nil)
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/login", LogIn).Methods("POST")

	for _, route := range routes {
		var handler http.Handler

		handler = tokenAuth.HandleFunc(route.HandlerFunc)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
