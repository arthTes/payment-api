package transaction

import operation "github.com/payment-api/internal/enum"

type Request struct {
	AccountID string         `json:"account_id" binding:"required"`
	Operation operation.Type `json:"operation_type" binding:"required"`
	Amount    float64        `json:"amount" binding:"required"`
}
