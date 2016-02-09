package main

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Username string        `json:"username" bson:"Username"`
	Password string        `json:"-" bson:"Password"`
	Token    string        `json:"token" bson:"-"`
}
