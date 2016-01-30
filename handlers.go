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
	col = "buildings"
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	session := context.Get(r, "db").(*mgo.Session)
	c := session.DB(database).C("users")

	result := User{}
	err := c.Find(bson.M{"Username": username}).One(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if password != result.Password {
		http.Error(w, "Invalid Login", http.StatusUnauthorized)
		return
	}

	t := memStore.NewToken(username)
	result.Token = t.Token

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func BuildList(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(col)
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

	c := session.DB(database).C(col)

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

func BuildInsert(w http.ResponseWriter, r *http.Request) {
	build := Building{}
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	build.Id = bson.NewObjectId()

	session := context.Get(r, "db").(*mgo.Session)

	c := session.DB(database).C(col)

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

	c := session.DB(database).C(col)

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

	c := session.DB(database).C(col)

	err := c.RemoveId(oid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
