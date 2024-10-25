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
		description   string
		input         domain.Transaction
		repository    repository.Transaction
		expectedError error
	}{
		{
			description: "success",
			input:       domain.Transaction{},
			repository: &transactionRepositoryMock{
				err: nil,
			},
			expectedError: nil,
		},
		{
			description: "any-persist-error",
			input:       domain.Transaction{},
			repository: &transactionRepositoryMock{
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

			TransactionUseCase := NewTransactionUseCase(scenario.repository)

			err := TransactionUseCase.Create(ctx, scenario.input)

			assert.Equal(t, scenario.expectedError, err)
		})
	}
}
