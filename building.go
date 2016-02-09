package main

import (
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

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

type Appartment struct {
	Title    string `json:"title" bson:"Title"`
	Position int    `json:"position" bson:"Position"`
}

type Building struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Address     Address       `json:"address" bson:"Address"`
	Oil         int64         `json:"oil" bson:"Oil"`
	Fund        int64         `json:"fund" bson:"Fund"`
	Active      bool          `json:"active" bson:"Active"`
	Managment   bool          `json:"managment" bson:"Managment"`
	Appartments []Appartment  `json:"appartments" bson:"Appartments"`
}

type BuildingExpense struct {
	Id         bson.ObjectId `json:"id" bson:"_id"`
	BuildId    string        `json:"buildId" bson:"BuildId"`
	Year       int           `json:"year" bson:"Year"`
	Month      int           `json:"month" bson:"Month"`
	Category   string        `json:"category" bson:"Category"`
	Percentage string        `json:"percentage" bson:"Percentage"`
	Amount     int64         `json:"amount" bson:"<Amount></Amount>"`
}

type PublicBuild Building

func (b Building) MarshalJSON() ([]byte, error) {
	if b.Appartments == nil {
		b.Appartments = []Appartment{}
	}
	return json.Marshal(struct {
		PublicBuild
		Title string `json:"title"`
	}{
		PublicBuild: PublicBuild(b),
		Title:       b.Address.Street + " " + b.Address.StreetNumber + ", " + b.Address.PostalCode + " " + b.Address.Area,
	})
}
