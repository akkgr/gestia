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

var memStore = New("salty")

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

	siteType := flag.String("type", "dir", "site path type zip or dir")
	zipPath := flag.String("path", "www", "path containing assets")
	flag.Parse()

	if *siteType == "zip" {
		rd, err := zip.OpenReader(*zipPath)
		if err != nil {
			log.Fatal(err)
		}
		fs := zipfs.New(rd, *zipPath)
		router.PathPrefix("/").Handler(http.FileServer(httpfs.New(fs)))
	} else {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(*zipPath)))
	}

	withdb := WithDB(db, router)

	withcors := corsHandler(withdb)

	withGz := gziphandler.GzipHandler(withcors)

	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", context.ClearHandler(withGz)))
}
