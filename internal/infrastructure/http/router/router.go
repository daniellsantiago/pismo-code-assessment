package router

import (
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler"
)

func New(accountHandler *handler.AccountHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /accounts", accountHandler.Create)

	return mux
}
