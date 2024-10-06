package di

import (
	"payment-integration/internal/gateway"
	"payment-integration/internal/repository"
	"payment-integration/internal/service"

	"go.mongodb.org/mongo-driver/mongo"
)

type Container struct {
	PaymentService *service.PaymentService
}

func NewContainer(db *mongo.Database) *Container {
	gatewayA := &gateway.GatewayA{}
	gatewayB := &gateway.GatewayB{}

	// Initialize the transaction repository
	transactionRepo := repository.NewTransactionRepository(db)

	// Pass the repositories to the payment service
	paymentService := service.NewPaymentService(gatewayA, gatewayB, transactionRepo)

	return &Container{
		PaymentService: paymentService,
	}
}
