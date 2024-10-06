package service_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"payment-integration/internal/model"
	"payment-integration/internal/repository"
	"payment-integration/internal/service"
)

type MockGatewayA struct{}
type MockGatewayB struct{}

func (g *MockGatewayA) Deposit(req model.TransactionRequest) (model.TransactionResponse, error) {
	return model.TransactionResponse{Status: "Success"}, nil
}

func (g *MockGatewayA) Withdraw(req model.TransactionRequest) (model.TransactionResponse, error) {
	return model.TransactionResponse{Status: "Success"}, nil
}

func (g *MockGatewayA) HandleCallback(data []byte) error {
	return nil
}

func (g *MockGatewayB) Deposit(req model.TransactionRequest) (model.TransactionResponse, error) {
	return model.TransactionResponse{Status: "Success"}, nil
}

func (g *MockGatewayB) Withdraw(req model.TransactionRequest) (model.TransactionResponse, error) {
	return model.TransactionResponse{Status: "Success"}, nil
}

func (g *MockGatewayB) HandleCallback(data []byte) error {
	return nil
}

func TestDepositHandler(t *testing.T) {
	mockRepo := &repository.TransactionRepository{
		// To DO: Mock MongoDB methods to execute and pass this test case
	}

	paymentService := service.NewPaymentService(&MockGatewayA{}, &MockGatewayB{}, mockRepo)

	reqBody := model.TransactionRequest{
		Amount:   100.0,
		Currency: "EUR",
		Gateway:  "A",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	paymentService.DepositHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestWithdrawHandler(t *testing.T) {
	mockRepo := &repository.TransactionRepository{
		// To DO: Mock MongoDB methods to execute and pass this test case
	}

	paymentService := service.NewPaymentService(&MockGatewayA{}, &MockGatewayB{}, mockRepo)

	reqBody := model.TransactionRequest{
		Amount:   100.0,
		Gateway:  "B",
		Currency: "EUR",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	paymentService.WithdrawHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
