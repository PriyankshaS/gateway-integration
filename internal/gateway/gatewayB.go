package gateway

import (
	"payment-integration/internal/model"
)

type GatewayB struct{}

func NewGatewayB() GatewayB {
	return GatewayB{}
}

func (g GatewayB) Deposit(req model.TransactionRequest) (model.TransactionResponse, error) {
	// Simulate a request to Gateway B
	response := model.TransactionResponse{Status: "success", Message: "Deposited to Gateway B"}
	return response, nil
}

func (g GatewayB) Withdraw(req model.TransactionRequest) (model.TransactionResponse, error) {
	// Simulate a request to Gateway B
	response := model.TransactionResponse{Status: "success", Message: "Withdraw from Gateway B"}
	return response, nil
}

func (g *GatewayB) HandleCallback(data []byte) error {
	// Callback handling
	return nil
}
