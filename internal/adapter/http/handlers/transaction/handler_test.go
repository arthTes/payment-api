package transaction

import (
	"bytes"
	"context"
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

type transactionUseCaseMock struct {
	Result domain.Transaction
	err    error
}

func (a transactionUseCaseMock) Create(context.Context, domain.Transaction) error {
	return a.err
}

func (a transactionUseCaseMock) Get(context.Context, string) (domain.Transaction, error) {
	return a.Result, a.err
}

func Test_transactionCreateHandler(t *testing.T) {
	scenarios := []struct {
		description    string
		input          []byte
		useCase        usecase.TransactionUseCase
		expectedStatus int
	}{
		{
			description: "success",
			input:       []byte(`{"account_id": "any-account-id", "operation_type": 1,"amount": 10.1}`),
			useCase: &transactionUseCaseMock{
				Result: domain.Transaction{},
				err:    nil,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			description: "persistence error",
			input:       []byte(`{"account_id": "any-account-id", "operation_type": 1,"amount": 10.1}`),
			useCase: &transactionUseCaseMock{
				Result: domain.Transaction{},
				err:    exceptions.PersistenceError,
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			description: "request body empty",
			input:       []byte(`{}`),
			useCase: &transactionUseCaseMock{
				Result: domain.Transaction{},
				err:    nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "request body with any wrong content",
			input:       []byte(`{"any-information": "whatever"}`),
			useCase: &transactionUseCaseMock{
				Result: domain.Transaction{},
				err:    nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "amount negative",
			input:       []byte(`{"account_id": "any-account-id", "operation_type": 1,"amount": -10.1}`),
			useCase: &transactionUseCaseMock{
				Result: domain.Transaction{},
				err:    nil,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			description: "invalid operation type",
			input:       []byte(`{"account_id": "any-account-id", "operation_type": 10,"amount": 10.1}`),
			useCase: &transactionUseCaseMock{
				Result: domain.Transaction{},
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
			SetTransactionRoutes(ctx, router, scenario.useCase)

			request, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(scenario.input))

			router.ServeHTTP(rr, request)

			assert.Equal(t, scenario.expectedStatus, rr.Code)
		})
	}
}
