package transaction

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"

	"github.com/payment-api/infrastructure/exceptions"
	"github.com/payment-api/infrastructure/logger"
	"github.com/payment-api/infrastructure/telemetry"
	"github.com/payment-api/internal/domain"
	"github.com/payment-api/internal/usecase"
)

func SetTransactionRoutes(ctx context.Context, r *gin.Engine, s usecase.TransactionUseCase) {
	r.POST("/api/v1/transactions", createTransaction(ctx, s))
}

func createTransaction(ctx context.Context, transactionUseCase usecase.TransactionUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := telemetry.Span(ctx, "http:handler:createTransaction", trace.SpanKindServer)
		defer span.End()

		var request Request

		if err := c.BindJSON(&request); err != nil {
			telemetry.ErrorSpan(span, err)
			logger.Error(logger.HTTPError, "cannot marshal body")

			c.JSON(http.StatusBadRequest, map[string]string{
				"message": "invalid parameters",
			})
			return
		}

		if !request.Operation.IsValid() {
			telemetry.ErrorSpan(span, exceptions.InvalidParameterError)
			logger.Error(logger.HTTPError, "invalid operation parameter")

			c.JSON(http.StatusBadRequest, map[string]string{
				"message": "invalid operation parameter",
				"reason":  exceptions.InvalidOperationTypeError.Error(),
			})
			return
		}

		if request.Amount < 0 {
			telemetry.ErrorSpan(span, exceptions.InvalidParameterError)
			logger.Error(logger.HTTPError, "invalid amount parameter")

			c.JSON(http.StatusBadRequest, map[string]string{
				"message": "invalid amount parameter",
				"reason":  exceptions.InvalidAmountError.Error(),
			})
			return
		}

		transaction := domain.NewTransaction(request.AccountID, request.Operation, request.Amount)

		err := transactionUseCase.Create(ctx, transaction)
		if err != nil {
			telemetry.ErrorSpan(span, err)
			logger.Error(logger.HTTPError, "Error creating transaction")
			c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "failed create transaction",
				"reason":  err.Error(),
			})

			return
		}

		logger.Info(logger.ServerInfo, fmt.Sprintf("Transaction Created %v", transaction))

		c.JSON(http.StatusCreated, map[string]string{"success": "created"})
	}
}
