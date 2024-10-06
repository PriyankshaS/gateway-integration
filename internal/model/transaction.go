package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Transaction struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Amount   float64            `bson:"amount"`
	Currency string             `bson:"currency"`
	Status   string             `bson:"status"`
	Gateway  string             `bson:"gateway"`
}

type TransactionRequest struct {
	Amount   float64 `bson:"amount"`
	Currency string  `bson:"currency"`
	Gateway  string  `bson:"gateway"`
}
