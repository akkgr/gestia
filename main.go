package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
)

const (
	server   = "localhost"
	database = "estia"
)

func main() {
	log.SetOutput(os.Stderr)
	log.Printf(
		"%s\t%s",
		"Server listening on ",
		":8080",
	)

	db, err := mgo.Dial(server)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	router := NewRouter()

	withdb := WithDB(db, router)

	withcors := corsHandler(withdb)

	withGz := gziphandler.GzipHandler(withcors)

	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(withGz)))
}
