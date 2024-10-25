package account

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"github.com/payment-api/infrastructure/exceptions"
	"github.com/payment-api/infrastructure/logger"
	"github.com/payment-api/infrastructure/telemetry"
	"github.com/payment-api/internal/domain"
	"github.com/payment-api/internal/usecase"
)

func SetAccountRoutes(ctx context.Context, r *gin.Engine, s usecase.AccountUseCase) {
	r.POST("/api/v1/accounts", createAccount(ctx, s))
	r.GET("/api/v1/accounts/:account_id", getAccount(ctx, s))
}

func getAccount(ctx context.Context, accountUseCase usecase.AccountUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := telemetry.Span(ctx, "http:handler:getAccount", trace.SpanKindServer)
		defer span.End()

		accountID := c.Param("account_id")
		if accountID == "" {
			telemetry.ErrorSpan(span, exceptions.InvalidParameterError)
			logger.Error(logger.HTTPError, "missing account ID")

			c.JSON(http.StatusBadRequest, map[string]string{
				"message": "invalid parameters",
			})
			return
		}

		persistedAccount, err := accountUseCase.Get(ctx, accountID)
		if err != nil {
			telemetry.ErrorSpan(span, err)
			logger.Error(logger.HTTPError, "ErrorSpan creating account")
			c.JSON(http.StatusNotFound, map[string]string{
				"message": "failed get account",
				"reason":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, persistedAccount)
	}
}

func createAccount(ctx context.Context, accountUseCase usecase.AccountUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := telemetry.Span(ctx, "http:handler:createAccount", trace.SpanKindServer)
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

		generatedAccountID := uuid.New().String()
		account := domain.NewAccount(generatedAccountID, request.DocumentNumber)

		err := accountUseCase.Create(ctx, account)
		if err != nil {
			telemetry.ErrorSpan(span, err)
			logger.Error(logger.HTTPError, "ErrorSpan creating account")
			c.JSON(http.StatusUnprocessableEntity, map[string]string{
				"message": "failed create account",
				"reason":  err.Error(),
			})

			return
		}

		logger.Info(logger.ServerInfo, fmt.Sprintf("Account Created %v", account))

		c.JSON(http.StatusCreated, map[string]string{"success": "created", "id": generatedAccountID})
	}
}
