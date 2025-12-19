package main

import (
	"context"
	"errors"
	"log"
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
	"github.com/nubank/pismo-code-assessment/internal/usecase/account"
	"github.com/nubank/pismo-code-assessment/internal/usecase/transaction"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	accountRepo := database.NewAccountRepository(db)
	transactionRepo := database.NewTransactionRepository(db)

	// Account use cases and handler
	createAccount := account.NewCreateAccount(accountRepo)
	getAccount := account.NewGetAccount(accountRepo)
	accountHandler := handler.NewAccountHandler(createAccount, getAccount)

	// Transaction use cases and handler
	createTransaction := transaction.NewCreateTransaction(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(createTransaction)

	r := router.New(accountHandler, transactionHandler)

	srv := server.New(cfg.Server.Port, r)

	go func() {
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
