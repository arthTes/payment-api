package repository

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/payment-api/infrastructure/logger"
	"github.com/payment-api/infrastructure/postgres"
	"github.com/payment-api/infrastructure/telemetry"
	"github.com/payment-api/internal/domain"
)

type Transaction interface {
	Push(ctx context.Context, entity domain.Transaction) error
}

type transactionImpl struct {
	repository postgres.Repository
}

func (t transactionImpl) Push(ctx context.Context, entity domain.Transaction) error {
	ctx, span := telemetry.Span(ctx, "repository:transaction:Push", trace.SpanKindInternal)
	defer span.End()

	q := `
	INSERT INTO transactions (account_id, operation_type_id, amount, event_date)
        VALUES ($1, $2, $3, $4)
        RETURNING id;
    `

	err := t.repository.Push(q, entity.AccountID, entity.OperationType, entity.Amount, time.Now())
	if err != nil {
		telemetry.ErrorSpan(span, err)
		logger.Error(logger.ServerError, "Error pushing transaction to postgres", err)
		return err
	}

	return nil
}

func NewTransactionRepository(repository postgres.Repository) Transaction {
	return transactionImpl{repository: repository}
}
