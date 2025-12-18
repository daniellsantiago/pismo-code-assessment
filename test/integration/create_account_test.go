package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nubank/pismo-code-assessment/internal/infrastructure/http/dto"
)

func TestCreateAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	ts := SetupTestServer(t, ctx)

	t.Run("creates account successfully", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"document_number": "12345678900"}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/accounts", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response dto.CreateAccountResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.NotZero(t, response.AccountID)
		assert.Equal(t, "12345678900", response.DocumentNumber)
	})

	t.Run("returns 422 when document number already exists", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"document_number": "99999999999"}`)

		// when
		resp1, err := http.Post(ts.Server.URL+"/accounts", "application/json", body)
		require.NoError(t, err)
		resp1.Body.Close()

		// then
		assert.Equal(t, http.StatusCreated, resp1.StatusCode)

		body = bytes.NewBufferString(`{"document_number": "99999999999"}`)
		resp2, err := http.Post(ts.Server.URL+"/accounts", "application/json", body)
		require.NoError(t, err)
		defer resp2.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, resp2.StatusCode)
	})

	t.Run("returns 422 when document number is empty", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`{"document_number": ""}`)

		// when
		resp, err := http.Post(ts.Server.URL+"/accounts", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("returns 400 when body is invalid", func(t *testing.T) {
		// given
		body := bytes.NewBufferString(`invalid json`)

		// when
		resp, err := http.Post(ts.Server.URL+"/accounts", "application/json", body)
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
