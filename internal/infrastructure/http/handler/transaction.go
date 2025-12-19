package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/response"
)

//go:generate mockgen -source=transaction.go -destination=mocks/transaction_mock.go -package=mocks
type transactionCreator interface {
	Execute(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*domain.Transaction, error)
}

type TransactionHandler struct {
	createTransaction transactionCreator
	logger            *slog.Logger
}

func NewTransactionHandler(createTransaction transactionCreator, logger *slog.Logger) *TransactionHandler {
	return &TransactionHandler{
		createTransaction: createTransaction,
		logger:            logger,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body",
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

	transaction, err := h.createTransaction.Execute(r.Context(), req.AccountID, req.OperationTypeID, req.Amount)
	if err != nil {
		h.logger.Error("failed to create transaction",
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
	})
}
