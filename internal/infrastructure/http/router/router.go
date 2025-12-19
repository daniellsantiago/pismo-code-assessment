package router

import (
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler"
)

func New(accountHandler *handler.AccountHandler, transactionHandler *handler.TransactionHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /accounts", accountHandler.Create)
	mux.HandleFunc("GET /accounts/{accountId}", accountHandler.Get)
	mux.HandleFunc("POST /transactions", transactionHandler.Create)

	return mux
}
