package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
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
}

func NewAccountHandler(createAccount accountCreator, getAccount accountGetter) *AccountHandler {
	return &AccountHandler{
		createAccount: createAccount,
		getAccount:    getAccount,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	account, err := h.createAccount.Execute(r.Context(), req.DocumentNumber)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusCreated, dto.CreateAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountIDStr := r.PathValue("accountId")
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid account id")
		return
	}

	account, err := h.getAccount.Execute(r.Context(), accountID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.respondJSON(w, http.StatusOK, dto.GetAccountResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	})
}

func (h *AccountHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		h.respondError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidDocumentNumber),
		errors.Is(err, domain.ErrAccountAlreadyExists):
		h.respondError(w, http.StatusUnprocessableEntity, err.Error())
	default:
		h.respondError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *AccountHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AccountHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
