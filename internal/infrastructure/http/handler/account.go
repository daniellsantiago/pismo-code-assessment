package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
)

//go:generate mockgen -source=account.go -destination=mocks/account_mock.go -package=mocks
type accountCreator interface {
	Execute(ctx context.Context, documentNumber string) (*domain.Account, error)
}

type AccountHandler struct {
	createAccount accountCreator
}

func NewAccountHandler(createAccount accountCreator) *AccountHandler {
	return &AccountHandler{createAccount: createAccount}
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

func (h *AccountHandler) handleError(w http.ResponseWriter, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		h.respondError(w, http.StatusUnprocessableEntity, domainErr.Error())
		return
	}

	h.respondError(w, http.StatusInternalServerError, "internal server error")
}

func (h *AccountHandler) respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AccountHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, map[string]string{"error": message})
}
