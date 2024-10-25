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

type Account interface {
	Push(ctx context.Context, entity domain.Account) error
	Get(ctx context.Context, id string) (domain.Account, error)
}

type (
	accountImpl struct {
		repository postgres.Repository
	}

	result struct {
		Id             string
		DocumentNumber string
		CreatedAt      time.Time
	}
)

func (a *accountImpl) Get(ctx context.Context, id string) (domain.Account, error) {
	ctx, span := telemetry.Span(ctx, "repository:account:Get", trace.SpanKindInternal)
	defer span.End()

	q := `
		SELECT * FROM accounts WHERE id = $1;
    `

	var resultPersisted result
	err := a.repository.GetById(q, id, &resultPersisted.Id, &resultPersisted.DocumentNumber, &resultPersisted.CreatedAt)

	if err != nil {
		logger.Error(logger.ServerError, "Error getting account to postgres", err)
		return domain.Account{}, err
	}

	return domain.Account{
		Id:             resultPersisted.Id,
		DocumentNumber: resultPersisted.DocumentNumber,
	}, nil
}

func (a *accountImpl) Push(ctx context.Context, entity domain.Account) error {
	ctx, span := telemetry.Span(ctx, "repository:account:Push", trace.SpanKindInternal)
	defer span.End()

	q := `
	INSERT INTO accounts (id, document_number, created_at)
        VALUES ($1, $2, $3)
        RETURNING id;
    `

	err := a.repository.Push(q, entity.Id, entity.DocumentNumber, time.Now())
	if err != nil {
		telemetry.ErrorSpan(span, err)
		logger.Error(logger.ServerError, "Error pushing account to postgres", err)
		return err
	}

	return nil
}

func NewAccountRepository(repository postgres.Repository) Account {
	return &accountImpl{
		repository: repository,
	}
}
