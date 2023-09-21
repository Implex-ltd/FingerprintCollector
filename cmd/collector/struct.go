package main

import "go.mongodb.org/mongo-driver/mongo"

var (
	collection      *mongo.Collection
	visitcollection *mongo.Collection
)

type FpPayload struct {
	Data struct {
		N  string `json:"n"`
		ID string `json:"id"`
	} `json:"data"`
}

type FpPayloadRaw struct {
	Data struct {
		N string      `json:"n"`
		ID interface{} `json:"id"`
	} `json:"data"`
}