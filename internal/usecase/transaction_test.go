package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/payment-api/infrastructure/exceptions"
	"github.com/payment-api/internal/adapter/repository"
	"github.com/payment-api/internal/domain"
)

type transactionRepositoryMock struct {
	Result domain.Transaction
	err    error
}

func (r *transactionRepositoryMock) Push(_ context.Context, _ domain.Transaction) error {
	return r.err
}

func Test_TransactionCreateUseCase(t *testing.T) {
	scenarios := []struct {
		description           string
		input                 domain.Transaction
		accountRepository     repository.Account
		transactionRepository repository.Transaction
		expectedError         error
	}{
		{
			description: "success",
			input:       domain.Transaction{},
			accountRepository: &accountRepositoryMock{
				err: nil,
			},
			transactionRepository: &transactionRepositoryMock{
				err: nil,
			},
			expectedError: nil,
		},
		{
			description: "account-not-found",
			input:       domain.Transaction{},
			accountRepository: &accountRepositoryMock{
				err: exceptions.EntityNotFoundError,
			},
			transactionRepository: &transactionRepositoryMock{
				err: errors.New("persist error"),
			},
			expectedError: exceptions.EntityNotFoundError,
		},
		{
			description: "any-persist-error",
			input:       domain.Transaction{},
			accountRepository: &accountRepositoryMock{
				err: nil,
			},
			transactionRepository: &transactionRepositoryMock{
				err: errors.New("persist error"),
			},
			expectedError: exceptions.PersistenceError,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "service-name", "payment-api")

			traceProvider := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
			traceProvider.Tracer(ctx.Value("service-name").(string))

			TransactionUseCase := NewTransactionUseCase(scenario.accountRepository, scenario.transactionRepository)

			err := TransactionUseCase.Create(ctx, scenario.input)

			assert.Equal(t, scenario.expectedError, err)
		})
	}
}
