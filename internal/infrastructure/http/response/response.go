package response

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/nubank/pismo-code-assessment/internal/domain"
)

var kindToStatus = map[domain.ErrorKind]int{
	domain.KindValidation: http.StatusUnprocessableEntity,
	domain.KindNotFound:   http.StatusNotFound,
}

func JSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, map[string]string{"error": message})
}

func HandleError(w http.ResponseWriter, err error) {
	var domainErr *domain.Error
	if errors.As(err, &domainErr) {
		status := kindToStatus[domainErr.Kind]
		if status == 0 {
			status = http.StatusInternalServerError
		}
		Error(w, status, domainErr.Message)
		return
	}
	Error(w, http.StatusInternalServerError, "internal server error")
}
