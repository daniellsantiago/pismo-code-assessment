package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/config"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/database"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/router"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/server"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/logger"
	"github.com/nubank/pismo-code-assessment/internal/usecase/account"
	"github.com/nubank/pismo-code-assessment/internal/usecase/transaction"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.Environment)

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Error("failed to connect to database", "error", err.Error())
		os.Exit(1)
	}
	defer db.Close()

	// Repositories
	accountRepo := database.NewAccountRepository(db)
	transactionRepo := database.NewTransactionRepository(db)

	// Account use cases and handler
	createAccount := account.NewCreateAccount(accountRepo)
	getAccount := account.NewGetAccount(accountRepo)
	accountHandler := handler.NewAccountHandler(createAccount, getAccount, log)

	// Transaction use cases and handler
	createTransaction := transaction.NewCreateTransaction(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(createTransaction, log)

	r := router.New(accountHandler, transactionHandler)

	srv := server.New(cfg.Server.Port, r)

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", "error", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", "error", err.Error())
		os.Exit(1)
	}
}
