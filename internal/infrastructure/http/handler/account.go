package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/response"
)

type accountCreator interface {
	Execute(ctx context.Context, documentNumber string) (*domain.Account, error)
}

type accountGetter interface {
	Execute(ctx context.Context, accountID int64) (*domain.Account, error)
}

//go:generate mockgen -source=account.go -destination=mocks/account_mock.go -package=mocks
type AccountHandler struct {
	createAccount accountCreator
	getAccount    accountGetter
	logger        *slog.Logger
}

func NewAccountHandler(createAccount accountCreator, getAccount accountGetter, logger *slog.Logger) *AccountHandler {
	return &AccountHandler{
		createAccount: createAccount,
		getAccount:    getAccount,
		logger:        logger,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body",
			slog.String("error", err.Error()),
		)
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DocumentNumber == "" {
		response.Error(w, http.StatusBadRequest, "document_number is required")
		return
	}

	account, err := h.createAccount.Execute(r.Context(), req.DocumentNumber)
	if err != nil {
		h.logger.Error("failed to create account",
			slog.String("document_number", req.DocumentNumber),
			slog.String("error", err.Error()),
		)
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusCreated, dto.CreateAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.PathValue("accountId")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.getAccount.Execute(r.Context(), accountID)
	if err != nil {
		h.logger.Error("failed to get account",
			slog.Int64("account_id", accountID),
			slog.String("error", err.Error()),
		)
		response.HandleError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, dto.GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}
