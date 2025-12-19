package handler

import (
	"context"
	"encoding/json"
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
}

func NewTransactionHandler(createTransaction transactionCreator) *TransactionHandler {
	return &TransactionHandler{createTransaction: createTransaction}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
