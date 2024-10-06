package service

import (
	"context"
	"encoding/json"
	"net/http"
	"payment-integration/internal/gateway"
	"payment-integration/internal/model"
	"payment-integration/internal/repository"
	"payment-integration/utils"
	"time"

	"github.com/sony/gobreaker"
)

type PaymentService struct {
	gatewayA   gateway.Gateway
	gatewayB   gateway.Gateway
	cbA        *gobreaker.CircuitBreaker
	cbB        *gobreaker.CircuitBreaker
	repository *repository.TransactionRepository
}

func NewPaymentService(gatewayA gateway.Gateway, gatewayB gateway.Gateway, repo *repository.TransactionRepository) *PaymentService {
	cbA := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "GatewayA",
		Timeout: 5 * time.Second,
	})
	cbB := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "GatewayB",
		Timeout: 5 * time.Second,
	})
	return &PaymentService{
		gatewayA:   gatewayA,
		gatewayB:   gatewayB,
		repository: repo,
		cbA:        cbA,
		cbB:        cbB,
	}
}

func (s *PaymentService) DepositHandler(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogTransaction(err.Error())
		utils.HandleError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var response model.TransactionResponse
	var result interface{}
	var err error

	if req.Gateway == "A" {
		result, err = s.cbA.Execute(func() (interface{}, error) {
			return s.callWithTimeout(func() (model.TransactionResponse, error) {
				return s.gatewayA.Deposit(req)
			})
		})

		if err != nil {
			// If Gateway A fails, try Gateway B with retries
			response, err = s.retry(func() (model.TransactionResponse, error) {
				return s.callWithTimeout(func() (model.TransactionResponse, error) {
					return s.gatewayB.Deposit(req)
				})
			}, 3) // Retry 3 times
		}
	} else if req.Gateway == "B" {
		result, err = s.cbB.Execute(func() (interface{}, error) {
			return s.callWithTimeout(func() (model.TransactionResponse, error) {
				return s.gatewayB.Deposit(req)
			})
		})

		if err != nil {
			// If Gateway B fails, try Gateway A with retries
			response, err = s.retry(func() (model.TransactionResponse, error) {
				return s.callWithTimeout(func() (model.TransactionResponse, error) {
					return s.gatewayA.Deposit(req)
				})
			}, 3) // Retry 3 times
		}
	} else {
		utils.HandleError(w, "Unsupported gateway", http.StatusBadRequest)
		return
	}
	// Type assertion to convert interface{} to models.TransactionResponse
	response, ok := result.(model.TransactionResponse)
	if !ok {
		utils.HandleError(w, "Unexpected error", http.StatusBadRequest)
		return
	}

	// Save transaction with initial status
	transaction := &model.Transaction{
		Amount:   req.Amount,
		Currency: req.Currency,
		Status:   "pending",
	}
	err = s.repository.SaveTransaction(transaction)
	if err != nil {
		utils.LogTransaction(err.Error())
		utils.HandleError(w, "Failed to save transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *PaymentService) GatewayACallback(w http.ResponseWriter, r *http.Request) {
	// Handle asynchronous callbacks from gateway A
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Callback A received"))
}

func (s *PaymentService) GatewayBCallback(w http.ResponseWriter, r *http.Request) {
	// Handle asynchronous callbacks from gateway B
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Callback B received"))
}

func (s *PaymentService) WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	var req model.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.HandleError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var response model.TransactionResponse
	var err error

	if req.Gateway == "A" {
		response, err = s.gatewayA.Withdraw(req)
	} else if req.Gateway == "B" {
		response, err = s.gatewayB.Withdraw(req)
	} else {
		http.Error(w, "Unsupported gateway", http.StatusBadRequest)
		return
	}

	// Save transaction status in MongoDB
	transaction := &model.Transaction{
		Amount:   req.Amount,
		Currency: req.Currency,
		Status:   "pending",
	}
	err = s.repository.SaveTransaction(transaction)
	if err != nil {
		utils.HandleError(w, "Error saving transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Call with timeout function
func (s *PaymentService) callWithTimeout(fn func() (model.TransactionResponse, error)) (model.TransactionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	resultChan := make(chan model.TransactionResponse)
	go func() {
		result, err := fn()
		if err != nil {
			resultChan <- model.TransactionResponse{}
		} else {
			resultChan <- result
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return model.TransactionResponse{}, context.DeadlineExceeded
	}
}

func (s *PaymentService) retry(fn func() (model.TransactionResponse, error), attempts int) (model.TransactionResponse, error) {
	var err error
	var result model.TransactionResponse

	for i := 0; i < attempts; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		time.Sleep(1 * time.Second) // Wait before retrying
	}
	return model.TransactionResponse{}, err
}
