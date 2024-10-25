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

type AccountUseCase interface {
	Create(context.Context, domain.Account) error
	Get(context.Context, string) (domain.Account, error)
}

type AccountUcImpl struct {
	accountRepository repository.Account
}

func (a *AccountUcImpl) Get(ctx context.Context, id string) (domain.Account, error) {
	ctx, span := telemetry.Span(ctx, "useCase:account:Get", trace.SpanKindInternal)
	defer span.End()

	persistedAccount, err := a.accountRepository.Get(ctx, id)
	if err != nil {
		telemetry.ErrorSpan(span, err)
		logger.Error(logger.ServerError, fmt.Sprintf("account not found error: %v", err.Error()))
		return domain.Account{}, exceptions.EntityNotFoundError
	}

	return persistedAccount, nil
}

func (a *AccountUcImpl) Create(ctx context.Context, account domain.Account) error {
	ctx, span := telemetry.Span(ctx, "useCase:account:Create", trace.SpanKindInternal)
	defer span.End()

	err := a.accountRepository.Push(ctx, account)
	if err != nil {
		telemetry.ErrorSpan(span, err)
		logger.Error(logger.ServerError, fmt.Sprintf("cannot create account error: %v", err.Error()))
		return exceptions.PersistenceError
	}

	return nil
}

func NewAccountUseCase(accountRepository repository.Account) AccountUseCase {
	return &AccountUcImpl{accountRepository: accountRepository}
}
