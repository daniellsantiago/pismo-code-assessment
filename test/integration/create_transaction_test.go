package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
)

func TestCreateTransaction_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	ts := SetupTestServer(t, ctx)

	// Create an account first
	accountBody := bytes.NewBufferString(`{"document_number": "12345678900"}`)
	accountResp, err := http.Post(ts.Server.URL+"/accounts", "application/json", accountBody)
	require.NoError(t, err)
	defer accountResp.Body.Close()

	var accountResponse dto.CreateAccountResponse
	err = json.NewDecoder(accountResp.Body).Decode(&accountResponse)
	require.NoError(t, err)

	accountID := accountResponse.AccountID

	t.Run("creates payment transaction successfully", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"account_id": ` + toString(accountID) + `, "operation_type_id": 4, "amount": 123.45}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/transactions", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreateTransactionResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotZero(t, response.TransactionID)
		assert.Equal(t, accountID, response.AccountID)
		assert.Equal(t, 4, response.OperationTypeID)
		assert.Equal(t, 123.45, response.Amount)
	})

	t.Run("creates purchase transaction with negative amount", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"account_id": ` + toString(accountID) + `, "operation_type_id": 1, "amount": 50.0}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/transactions", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreateTransactionResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, -50.0, response.Amount)
	})

	t.Run("returns 404 when account does not exist", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"account_id": 999999, "operation_type_id": 1, "amount": 50.0}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/transactions", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("returns 422 when operation type is invalid", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"account_id": ` + toString(accountID) + `, "operation_type_id": 99, "amount": 50.0}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/transactions", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("returns 422 when amount is zero", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"account_id": ` + toString(accountID) + `, "operation_type_id": 1, "amount": 0}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/transactions", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`invalid json`)

		// when
		resp, err := http.Post(ts.Server.URL+"/transactions", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func toString(id int64) string {
	return fmt.Sprintf("%d", id)
}
