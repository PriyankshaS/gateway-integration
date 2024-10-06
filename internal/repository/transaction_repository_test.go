package repository_test

import (
	"context"
	"testing"

	"payment-integration/internal/model"
	"payment-integration/internal/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestTransactionRepositorySaveTransaction(t *testing.T) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ := mongo.Connect(context.Background(), clientOptions)
	db := client.Database("test")

	repo := repository.NewTransactionRepository(db)
	transaction := &model.Transaction{
		Amount:   100.0,
		Currency: "EUR",
		Status:   "pending",
		Gateway:  "gateway_a",
	}

	err := repo.SaveTransaction(transaction)
	if err != nil {
		t.Errorf("Failed to save transaction: %v", err)
	}
}
