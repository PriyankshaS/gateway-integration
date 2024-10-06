package di_test

import (
	"context"
	"payment-integration/internal/di"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestNewContainer(t *testing.T) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		t.Fatal(err)
	}

	container := di.NewContainer(client.Database("payment_gateway"))
	if container.PaymentService == nil {
		t.Error("Expected PaymentService to be initialized")
	}
}
