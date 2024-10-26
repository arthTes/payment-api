package usecase

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"

	"github.com/payment-api/infrastructure/exceptions"
	"github.com/payment-api/infrastructure/logger"
	"github.com/payment-api/infrastructure/telemetry"
	"github.com/payment-api/internal/adapter/repository"
	"github.com/payment-api/internal/domain"
)

type (
	TransactionUseCase interface {
		Create(context.Context, domain.Transaction) error
	}
)

type (
	TransactionUcImpl struct {
		accountRepository     repository.Account
		transactionRepository repository.Transaction
	}
)

func (t TransactionUcImpl) Create(ctx context.Context, transaction domain.Transaction) error {
	ctx, span := telemetry.Span(ctx, "useCase:transaction:Create", trace.SpanKindInternal)
	defer span.End()

	_, err := t.accountRepository.Get(ctx, transaction.AccountID)
	if err != nil {
		telemetry.ErrorSpan(span, err)
		logger.Error(logger.ServerError, "account not found", err.Error())
		return exceptions.EntityNotFoundError
	}

	if err := t.transactionRepository.Push(ctx, transaction); err != nil {
		telemetry.ErrorSpan(span, err)
		logger.Error(logger.ServerError, fmt.Sprintf("cannot create transaction error: %v", err.Error()))
		return exceptions.PersistenceError
	}

	return nil
}

func NewTransactionUseCase(accountRepository repository.Account, transactionRepository repository.Transaction) TransactionUseCase {
	return TransactionUcImpl{
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}
