package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
