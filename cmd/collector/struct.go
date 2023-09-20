package main

import "go.mongodb.org/mongo-driver/mongo"

var (
	collection      *mongo.Collection
	visitcollection *mongo.Collection
)

type FpPayload struct {
	Data struct {
		N string `json:"n"`
		F string `json:"f"`
		D string `json:"d"`
		J string `json:"j"`
	} `json:"data"`
}

type FpPayloadRaw struct {
	Data struct {
		N string      `json:"n"`
		F interface{} `json:"f"`
		D interface{} `json:"d"`
		J interface{} `json:"j"`
	} `json:"data"`
}

type Fpjs struct {
	VisitorID string `json:"visitorId"`
}
