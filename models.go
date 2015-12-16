package main

import "gopkg.in/mgo.v2/bson"

type User struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Usename  string        `json:"username" bson:"Username"`
	Password string        `json:"password" bson:"Password"`
}

type GeoJson struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

type Address struct {
	Area         string  `json:"area" bson:"Area"`
	Street       string  `json:"street" bson:"Street"`
	StreetNumber string  `json:"streetnumber" bson:"StreetNumber"`
	PostalCode   string  `json:"postalcode" bson:"PostalCode"`
	Location     GeoJson `json:"location" bson:"Location"`
}

type Building struct {
	Id      bson.ObjectId `json:"id" bson:"_id"`
	Address Address       `json:"address" bson:"Address"`
}
