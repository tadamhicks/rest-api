package models

import "gopkg.in/mgo.v2/bson"

type Person struct {
	ID        bson.ObjectId        `json:"id" bson:"_id,omitempty"`
	Firstname string               `json:"firstname" bson:"firstname,omitempty"`
	Lastname  string               `json:"lastname" bson:"lastname,omitempty"`
	Address   *Address             `json:"address" bson:"address,omitempty"`
}

type Address struct {
	City    string  `json:"city" bson:"city,omitempty"`
	State   string  `json:"state" bson:"state,omitempty"`
}
