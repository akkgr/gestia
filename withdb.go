package main

import (
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
)

func WithDB(s *mgo.Session, h http.Handler) http.Handler {
	return &dbwrapper{dbSession: s, h: h}
}

type dbwrapper struct {
	dbSession *mgo.Session
	h         http.Handler
}

func (dbwrapper *dbwrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// copy the session
	dbcopy := dbwrapper.dbSession.Copy()
	defer dbcopy.Close() // clean up after

	// put the session in the context for this Request
	context.Set(r, "db", dbcopy)

	// serve the request
	dbwrapper.h.ServeHTTP(w, r)
}
