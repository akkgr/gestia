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
	log.Fatal(http.ListenAndServe(":8080", router))
}
