package domain

import "github.com/payment-api/internal/enum"

type Transaction struct {
	AccountID     string
	OperationType operation.Type
	Amount        float64
}

func NewTransaction(accountId string, operationType operation.Type, amount float64) Transaction {
	return Transaction{
		AccountID:     accountId,
		OperationType: operationType,
		Amount:        amount,
	}
}
