package gateway

import (
	"payment-integration/internal/model"
)

type Gateway interface {
	Deposit(req model.TransactionRequest) (model.TransactionResponse, error)
	Withdraw(req model.TransactionRequest) (model.TransactionResponse, error)
	HandleCallback(data []byte) error
}
