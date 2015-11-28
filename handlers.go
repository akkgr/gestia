package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func BuildIndex(w http.ResponseWriter, r *http.Request) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("estia").C("buildings")

	result := []Building{}
	err = c.Find(nil).All(&result)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

func BuildShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintln(w, "Todo show:", id)
}
