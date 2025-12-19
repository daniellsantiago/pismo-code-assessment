package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/handler/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTransactionHandler_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCreator := mocks.NewMocktransactionCreator(ctrl)
	handler := NewTransactionHandler(mockCreator)

	t.Run("creates transaction successfully", func(t *testing.T) {
		expectedTransaction := &domain.Transaction{
			ID:              1,
			AccountID:       1,
			OperationTypeID: domain.OperationTypePayment,
			Amount:          123.45,
			EventDate:       time.Now(),
		}

		mockCreator.EXPECT().
			Execute(gomock.Any(), int64(1), 4, 123.45).
			Return(expectedTransaction, nil)

		body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 4, "amount": 123.45}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response dto.CreateTransactionResponse
		json.NewDecoder(rec.Body).Decode(&response)

		assert.Equal(t, int64(1), response.TransactionID)
		assert.Equal(t, int64(1), response.AccountID)
		assert.Equal(t, 4, response.OperationTypeID)
		assert.Equal(t, 123.45, response.Amount)
	})

	t.Run("returns bad request when body is invalid", func(t *testing.T) {
		body := bytes.NewBufferString(`invalid json`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("returns not found when account does not exist", func(t *testing.T) {
		mockCreator.EXPECT().
			Execute(gomock.Any(), int64(999), 1, 50.0).
			Return(nil, domain.ErrAccountNotFound)

		body := bytes.NewBufferString(`{"account_id": 999, "operation_type_id": 1, "amount": 50.0}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns unprocessable entity when operation type is invalid", func(t *testing.T) {
		mockCreator.EXPECT().
			Execute(gomock.Any(), int64(1), 99, 50.0).
			Return(nil, domain.ErrInvalidOperationType)

		body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 99, "amount": 50.0}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("returns unprocessable entity when amount is invalid", func(t *testing.T) {
		mockCreator.EXPECT().
			Execute(gomock.Any(), int64(1), 1, 0.0).
			Return(nil, domain.ErrInvalidAmount)

		body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 1, "amount": 0}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("returns internal server error when error is unknown", func(t *testing.T) {
		mockCreator.EXPECT().
			Execute(gomock.Any(), int64(1), 1, 50.0).
			Return(nil, errors.New("database error"))

		body := bytes.NewBufferString(`{"account_id": 1, "operation_type_id": 1, "amount": 50.0}`)
		req := httptest.NewRequest(http.MethodPost, "/transactions", body)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
