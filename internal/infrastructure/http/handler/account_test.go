package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAccountHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreator := mocks.NewMockaccountCreator(ctrl)
	handler := NewAccountHandler(mockCreator)

	t.Run("creates account successfully", func(t *testing.T) {
		expectedAccount := &domain.Account{
			ID:             1,
			DocumentNumber: "12345678900",
		}

		mockCreator.EXPECT().
			Execute(gomock.Any(), "12345678900").
			Return(expectedAccount, nil)

		body := bytes.NewBufferString(`{"document_number": "12345678900"}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.CreateAccountResponse
		json.NewDecoder(rec.Body).Decode(&response)

		assert.Equal(t, int64(1), response.AccountID)
		assert.Equal(t, "12345678900", response.DocumentNumber)
	})

	t.Run("returns bad request when body is invalid", func(t *testing.T) {
		body := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns unprocessable entity when document number is empty", func(t *testing.T) {
		mockCreator.EXPECT().
			Execute(gomock.Any(), "").
			Return(nil, domain.ErrInvalidDocumentNumber)

		body := bytes.NewBufferString(`{"document_number": ""}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("returns unprocessable entity when account already exists", func(t *testing.T) {
		mockCreator.EXPECT().
			Execute(gomock.Any(), "99999999999").
			Return(nil, domain.ErrAccountAlreadyExists)

		body := bytes.NewBufferString(`{"document_number": "99999999999"}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})
}
