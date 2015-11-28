package main

import "gopkg.in/mgo.v2/bson"

type Address struct {
	Area         string `json:"area"`
	Street       string `json:"street"`
	StreetNumber string `json:"streetNumber"`
	PostalCode   string `json:"postalCode"`
}

type Building struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Address Address       `json:"address"`
}
