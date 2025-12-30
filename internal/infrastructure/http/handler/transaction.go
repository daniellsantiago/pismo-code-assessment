package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/response"
	"github.com/nubank/pismo-code-assessment/pkg/logger"
)

//go:generate mockgen -source=transaction.go -destination=mocks/transaction_mock.go -package=mocks
type transactionCreator interface {
	Execute(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*domain.Transaction, error)
}

type TransactionHandler struct {
	createTransaction transactionCreator
}

func NewTransactionHandler(createTransaction transactionCreator) *TransactionHandler {
	return &TransactionHandler{
		createTransaction: createTransaction,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(ctx, "failed to decode request body",
			slog.String("error", err.Error()),
		)
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.AccountID == 0 {
		response.Error(w, http.StatusBadRequest, "account_id is required")
		return
	}
	if req.OperationTypeID == 0 {
		response.Error(w, http.StatusBadRequest, "operation_type_id is required")
		return
	}
	if req.Amount == 0 {
		response.Error(w, http.StatusBadRequest, "amount is required")
		return
	}

	transaction, err := h.createTransaction.Execute(ctx, req.AccountID, req.OperationTypeID, req.Amount)
	if err != nil {
		logger.Error(ctx, "failed to create transaction",
			slog.Int64("account_id", req.AccountID),
			slog.Int("operation_type_id", req.OperationTypeID),
			slog.Float64("amount", req.Amount),
			slog.String("error", err.Error()),
		)
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.CreateTransactionResponse{
		TransactionID:   transaction.ID,
		AccountID:       transaction.AccountID,
		OperationTypeID: int(transaction.OperationTypeID),
		Amount:          transaction.Amount,
		Balance:         transaction.Balance,
	})
}
