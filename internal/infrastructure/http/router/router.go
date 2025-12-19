package router

import (
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/middleware"
)

func New(
	accountHandler *handler.AccountHandler,
	transactionHandler *handler.TransactionHandler,
	healthHandler *handler.HealthHandler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler.Check)
	mux.HandleFunc("POST /accounts", accountHandler.Create)
	mux.HandleFunc("GET /accounts/{accountId}", accountHandler.Get)
	mux.HandleFunc("POST /transactions", transactionHandler.Create)

	return middleware.Chain(
		mux,
		middleware.RequestID,
		middleware.Recoverer,
		middleware.Logger,
	)
}
