package gateway_test

import (
	"payment-integration/internal/gateway"
	"payment-integration/internal/model"
	"testing"
)

func TestGatewayADeposit(t *testing.T) {
	g := &gateway.GatewayA{}
	transactionRes, err := g.Deposit(model.TransactionRequest{Amount: 100.0, Currency: "EUR"})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if transactionRes.Status == "" {
		t.Error("Expected a transaction status, got empty string")
	}
}
