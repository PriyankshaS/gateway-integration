package main

import (
	"context"
	"log"
	"net/http"

	"payment-integration/internal/di"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Set up MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Ping the MongoDB server
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	// Initialize DI Container
	container := di.NewContainer(client.Database("payment_gateway"))

	router := mux.NewRouter()

	router.HandleFunc("/deposit", container.PaymentService.DepositHandler).Methods("POST")
	router.HandleFunc("/withdraw", container.PaymentService.WithdrawHandler).Methods("POST")
	router.HandleFunc("/callback/gateway_a", container.PaymentService.GatewayACallback).Methods("POST")
	router.HandleFunc("/callback/gateway_b", container.PaymentService.GatewayBCallback).Methods("POST")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
