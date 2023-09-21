package main

import "go.mongodb.org/mongo-driver/mongo"

var (
	collection      *mongo.Collection
	visitcollection *mongo.Collection
)

type FpPayload struct {
	N  string `json:"n"`
		ID string `json:"id"`
}
