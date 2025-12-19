package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestAccountHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreator := mocks.NewMockaccountCreator(ctrl)
	mockGetter := mocks.NewMockaccountGetter(ctrl)
	handler := NewAccountHandler(mockCreator, mockGetter, newTestLogger())

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

	t.Run("returns bad request when document number is empty", func(t *testing.T) {
		body := bytes.NewBufferString(`{"document_number": ""}`)
		req := httptest.NewRequest(http.MethodPost, "/accounts", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
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

func TestAccountHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreator := mocks.NewMockaccountCreator(ctrl)
	mockGetter := mocks.NewMockaccountGetter(ctrl)
	handler := NewAccountHandler(mockCreator, mockGetter, newTestLogger())

	t.Run("retrieves account successfully", func(t *testing.T) {
		expectedAccount := &domain.Account{
			ID:             1,
			DocumentNumber: "12345678900",
		}

		mockGetter.EXPECT().
			Execute(gomock.Any(), int64(1)).
			Return(expectedAccount, nil)

		req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
		req.SetPathValue("accountId", "1")
		rec := httptest.NewRecorder()

		handler.Get(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response dto.GetAccountResponse
		json.NewDecoder(rec.Body).Decode(&response)

		assert.Equal(t, int64(1), response.AccountID)
		assert.Equal(t, "12345678900", response.DocumentNumber)
	})

	t.Run("returns bad request when account id is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/accounts/invalid", nil)
		req.SetPathValue("accountId", "invalid")
		rec := httptest.NewRecorder()

		handler.Get(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns bad request when account id is blank", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/accounts/", nil)
		req.SetPathValue("accountId", "")
		rec := httptest.NewRecorder()

		handler.Get(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns not found when account does not exist", func(t *testing.T) {
		mockGetter.EXPECT().
			Execute(gomock.Any(), int64(999)).
			Return(nil, domain.ErrAccountNotFound)

		req := httptest.NewRequest(http.MethodGet, "/accounts/999", nil)
		req.SetPathValue("accountId", "999")
		rec := httptest.NewRecorder()

		handler.Get(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns internal server error when error is unknown", func(t *testing.T) {
		mockGetter.EXPECT().
			Execute(gomock.Any(), int64(1)).
			Return(nil, errors.New("database error"))

		req := httptest.NewRequest(http.MethodGet, "/accounts/1", nil)
		req.SetPathValue("accountId", "1")
		rec := httptest.NewRecorder()

		handler.Get(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
