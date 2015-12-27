package main

import (
	"log"
	"net/http"
	"os"
)

const (
	server   = "localhost"
	database = "estia"
)

func main() {

	router := NewRouter()
	log.SetOutput(os.Stderr)
	log.Printf(
		"%s\t%s",
		"Server listening on ",
		":8080",
	)
	log.Fatal(http.ListenAndServe(":8080", corsHandler(router)))
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusOK)
		} else {
			h.ServeHTTP(w, r)
		}
	}
}
