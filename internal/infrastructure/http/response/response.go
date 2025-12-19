package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/domain"
)

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidDocumentNumber),
		errors.Is(err, domain.ErrAccountAlreadyExists),
		errors.Is(err, domain.ErrInvalidOperationType),
		errors.Is(err, domain.ErrInvalidAmount):
		Error(w, http.StatusUnprocessableEntity, err.Error())
	default:
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}
