package repository

import (
	"context"

	"payment-integration/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TransactionRepository struct {
	collection *mongo.Collection
}

// NewTransactionRepository creates a new TransactionRepository
func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	return &TransactionRepository{
		collection: db.Collection("transactions"),
	}
}

// SaveTransaction saves a transaction to the database
func (r *TransactionRepository) SaveTransaction(transaction *model.Transaction) error {
	_, err := r.collection.InsertOne(context.Background(), transaction)
	return err
}

// UpdateTransactionStatus updates the status of a transaction
func (r *TransactionRepository) UpdateTransactionStatus(id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"status": status}})
	return err
}
