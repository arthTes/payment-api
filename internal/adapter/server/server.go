package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/payment-api/config"
	"github.com/payment-api/infrastructure/logger"
	"github.com/payment-api/infrastructure/postgres"
	"github.com/payment-api/internal/adapter/http/handlers/account"
	"github.com/payment-api/internal/adapter/http/handlers/transaction"
	"github.com/payment-api/internal/adapter/http/middlewares"
	"github.com/payment-api/internal/adapter/repository"
	"github.com/payment-api/internal/usecase"
)

type Server struct {
	config   config.Configuration
	services svs
}

type svs struct {
	account     usecase.AccountUseCase
	transaction usecase.TransactionUseCase
}

func New(ctx context.Context, cfg config.Configuration) (a Server) {
	a.config = cfg

	pgRepository, err := postgres.NewRepository(ctx, a.config.Postgres.Url)
	if err != nil {
		logger.Fatal(logger.ConfigError, fmt.Sprintf("Cannot connect postgresql error: %v", err))
	}

	accountRepository := repository.NewAccountRepository(*pgRepository)
	a.services.account = usecase.NewAccountUseCase(accountRepository)

	transactionRepository := repository.NewTransactionRepository(*pgRepository)
	a.services.transaction = usecase.NewTransactionUseCase(transactionRepository)

	return a
}

func (a *Server) Run(ctx context.Context, cancel context.CancelFunc) func() error {
	return func() error {
		defer cancel()

		router := gin.Default()

		router.Use(middlewares.Recover())

		account.SetAccountRoutes(ctx, router, a.services.account)
		transaction.SetTransactionRoutes(ctx, router, a.services.transaction)

		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", a.config.Server.Port),
			Handler: router,
		}

		go shutdown(ctx, server)
		return server.ListenAndServe()
	}
}

func shutdown(ctx context.Context, server *http.Server) {
	<-ctx.Done()
	logger.Info(logger.ServerInfo, "New we can do an shutdown gracefully")
	server.Shutdown(ctx)
}
