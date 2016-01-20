package main

import (
	"archive/zip"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/context"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
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

	zipPath := flag.String("zip", "estia.static", "zip file containing assets")
	flag.Parse()

	rd, err := zip.OpenReader(*zipPath)
	if err != nil {
		log.Fatal(err)
	}
	fs := zipfs.New(rd, *zipPath)
	router.PathPrefix("/").Handler(http.FileServer(httpfs.New(fs)))

	withdb := WithDB(db, router)

	withcors := corsHandler(withdb)

	withGz := gziphandler.GzipHandler(withcors)

	log.Fatal(http.ListenAndServe(":8080", context.ClearHandler(withGz)))
}
