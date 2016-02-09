package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	buildCollection = "buildings"
)

func BuildList(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(buildCollection)
	result := []Building{}

	err := c.Find(nil).Sort("Address.Street", "Address.StreetNumber", "Address.Area").All(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func BuildById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(buildCollection)

	result := Building{}
	err := c.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func BuildByCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inact := vars["inact"]
	act := vars["act"]
	mng := vars["mng"]
	nomng := vars["nomng"]

	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(buildCollection)

	result := []Building{}
	filter := bson.M{}

	if inact == "true" && act == "false" {
		filter["Active"] = false
	}

	if inact == "false" && act == "true" {
		filter["Active"] = true
	}

	if nomng == "true" && mng == "false" {
		filter["Managment"] = false
	}

	if nomng == "false" && mng == "true" {
		filter["Managment"] = true
	}

	if inact == "true" && act == "true" && nomng == "true" && mng == "true" {
		filter = nil
	}

	if inact == "false" && act == "false" && nomng == "false" && mng == "false" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
		return
	}

	err := c.Find(filter).Sort("Address.Street", "Address.StreetNumber", "Address.Area").All(&result)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func BuildInsert(w http.ResponseWriter, r *http.Request) {
	build := Building{}
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	build.Id = bson.NewObjectId()

	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(buildCollection)

	err := c.Insert(build)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(build)
}

func BuildUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	build := Building{}
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(buildCollection)

	err := c.UpdateId(oid, build)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(build)
}

func BuildDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(id)

	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(buildCollection)

	err := c.RemoveId(oid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
