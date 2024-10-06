package gateway

import (
	"payment-integration/internal/model"
)

type GatewayA struct{}

func NewGatewayA() GatewayA {
	return GatewayA{}
}

func (g GatewayA) Deposit(req model.TransactionRequest) (model.TransactionResponse, error) {
	// Simulate a request to Gateway A
	response := model.TransactionResponse{Status: "success", Message: "Deposited to Gateway A"}
	return response, nil
}

func (g GatewayA) Withdraw(req model.TransactionRequest) (model.TransactionResponse, error) {
	// Simulate a request to Gateway A
	response := model.TransactionResponse{Status: "success", Message: "Withdraw from Gateway A"}
	return response, nil
}

func (g *GatewayA) HandleCallback(data []byte) error {
	// Callback handling
	return nil
}
