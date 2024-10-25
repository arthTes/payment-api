package account

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/payment-api/infrastructure/exceptions"
	"github.com/payment-api/internal/domain"
	"github.com/payment-api/internal/usecase"
)

type accountUseCaseMock struct {
	Result domain.Account
	err    error
}

func (a accountUseCaseMock) Create(context.Context, domain.Account) error {
	return a.err
}

func (a accountUseCaseMock) Get(context.Context, string) (domain.Account, error) {
	return a.Result, a.err
}

func Test_AccountCreateHandler(t *testing.T) {
	scenarios := []struct {
		description    string
		input          []byte
		useCase        usecase.AccountUseCase
		expectedStatus int
	}{
		{
			description: "success",
			input:       []byte(`{"document_number":"any-document"}`),
			useCase: &accountUseCaseMock{
				Result: domain.Account{
					Id:             "generated-account-id",
					DocumentNumber: "any-document",
				},
				err: nil,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			description: "persistence error",
			input:       []byte(`{"document_number":"any-document"}`),
			useCase: &accountUseCaseMock{
				Result: domain.Account{},
				err:    exceptions.PersistenceError,
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			description: "request body empty",
			input:       []byte(`{}`),
			useCase: &accountUseCaseMock{
				Result: domain.Account{},
				err:    nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "request body without document number",
			input:       []byte(`{"any-information": "whatever"}`),
			useCase: &accountUseCaseMock{
				Result: domain.Account{},
				err:    nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "service-name", "payment-api")

			traceProvider := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
			traceProvider.Tracer(ctx.Value("service-name").(string))

			rr := httptest.NewRecorder()
			router := gin.Default()
			SetAccountRoutes(ctx, router, scenario.useCase)

			request, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(scenario.input))

			router.ServeHTTP(rr, request)

			assert.Equal(t, scenario.expectedStatus, rr.Code)
		})
	}
}

func Test_AccountGetHandler(t *testing.T) {
	scenarios := []struct {
		description    string
		input          string
		useCase        usecase.AccountUseCase
		expectedStatus int
	}{
		{
			description: "success",
			input:       "any-valid-account-id",
			useCase: &accountUseCaseMock{
				Result: domain.Account{
					Id:             "generated-account-id",
					DocumentNumber: "any-document",
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
		},
		{
			description: "persistence error",
			input:       "any-valid-account-id",
			useCase: &accountUseCaseMock{
				Result: domain.Account{},
				err:    exceptions.PersistenceError,
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			description: "empty account id parameter",
			input:       "",
			useCase: &accountUseCaseMock{
				Result: domain.Account{},
				err:    nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.description, func(t *testing.T) {
			ctx := context.Background()
			ctx = context.WithValue(ctx, "service-name", "payment-api")

			traceProvider := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
			traceProvider.Tracer(ctx.Value("service-name").(string))

			rr := httptest.NewRecorder()
			router := gin.Default()
			SetAccountRoutes(ctx, router, scenario.useCase)

			request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", scenario.input), nil)

			router.ServeHTTP(rr, request)

			assert.Equal(t, scenario.expectedStatus, rr.Code)
		})
	}
}
