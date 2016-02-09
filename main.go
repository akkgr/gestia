package main

import (
	"archive/zip"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
	"gopkg.in/mgo.v2"
)

var memStore = NewMemoryTokenStore("salty")
var server string
var database string

func main() {
	log.SetOutput(os.Stderr)

	port := flag.String("port", "8080", "server listening tcp port")
	dbServer := flag.String("server", "localhost", "database server")
	dbName := flag.String("db", "estia", "database name")
	siteType := flag.String("type", "dir", "site path type zip or dir")
	sitePath := flag.String("path", "wwwroot", "path containing site")
	flag.Parse()

	server = *dbServer
	database = *dbName

	db, err := mgo.Dial(server)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	router := NewRouter()
	if *siteType == "zip" {
		rd, err := zip.OpenReader(*sitePath)
		if err != nil {
			log.Fatal(err)
		}
		fs := zipfs.New(rd, *sitePath)
		router.PathPrefix("/").Handler(http.FileServer(httpfs.New(fs)))
	} else {
		router.PathPrefix("/").Handler(http.FileServer(http.Dir(*sitePath)))
	}

	withLog := handlers.LoggingHandler(os.Stdout, router)

	withdb := WithDB(db, withLog)

	withcors := handlers.CORS()(withdb)

	withGz := handlers.CompressHandler(withcors)

	log.Printf(
		"%s\t%s",
		"Server listening on ",
		*port,
	)

	log.Fatal(http.ListenAndServeTLS(":"+*port, "cert.pem", "key.pem", context.ClearHandler(withGz)))
}
