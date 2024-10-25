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

type accountRepositoryMock struct {
	Result domain.Account
	err    error
}

func (r *accountRepositoryMock) Get(_ context.Context, _ string) (domain.Account, error) {
	return r.Result, r.err
}

func (r *accountRepositoryMock) Push(_ context.Context, _ domain.Account) error {
	return r.err
}

func Test_AccountCreateUseCase(t *testing.T) {
	scenarios := []struct {
		description   string
		input         domain.Account
		repository    repository.Account
		expectedError error
	}{
		{
			description: "success",
			input: domain.Account{
				Id:             "generated-account-id",
				DocumentNumber: "any-document",
			},
			repository: &accountRepositoryMock{
				Result: domain.Account{
					Id:             "generated-account-id",
					DocumentNumber: "any-document",
				},
				err: nil,
			},
			expectedError: nil,
		},
		{
			description: "any-persist-error",
			input: domain.Account{
				Id:             "generated-account-id",
				DocumentNumber: "any-document",
			},
			repository: &accountRepositoryMock{
				Result: domain.Account{},
				err:    errors.New("any-error"),
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

			accountUseCase := NewAccountUseCase(scenario.repository)

			err := accountUseCase.Create(ctx, scenario.input)

			assert.Equal(t, scenario.expectedError, err)
		})
	}
}

func Test_AccountGetUseCase(t *testing.T) {
	scenarios := []struct {
		description    string
		input          string
		repository     repository.Account
		expectedOutput domain.Account
		expectedError  error
	}{
		{
			description: "success",
			input:       "generated-account-id",
			repository: &accountRepositoryMock{
				Result: domain.Account{
					Id:             "generated-account-id",
					DocumentNumber: "any-document",
				},
				err: nil,
			},
			expectedOutput: domain.Account{
				Id:             "generated-account-id",
				DocumentNumber: "any-document",
			},
			expectedError: nil,
		},
		{
			description: "not-found-error",
			input:       "generated-account-id",
			repository: &accountRepositoryMock{
				Result: domain.Account{},
				err:    errors.New("any-error"),
			},
			expectedOutput: domain.Account{},
			expectedError:  exceptions.EntityNotFoundError,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "service-name", "payment-api")

			traceProvider := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
			traceProvider.Tracer(ctx.Value("service-name").(string))

			accountUseCase := NewAccountUseCase(scenario.repository)

			output, err := accountUseCase.Get(ctx, scenario.input)

			assert.Equal(t, scenario.expectedOutput, output)
			assert.Equal(t, scenario.expectedError, err)
		})
	}
}
