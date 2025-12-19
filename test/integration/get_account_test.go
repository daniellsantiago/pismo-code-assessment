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

func TestGetAccount_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	ts := SetupTestServer(t, ctx)

	t.Run("retrieves account successfully", func(t *testing.T) {
		// given - create an account first
		body := bytes.NewBufferString(`{"document_number": "11111111111"}`)
		createResp, err := http.Post(ts.Server.URL+"/accounts", "application/json", body)
		require.NoError(t, err)
		defer createResp.Body.Close()

		var createResponse dto.CreateAccountResponse
		err = json.NewDecoder(createResp.Body).Decode(&createResponse)
		require.NoError(t, err)

		// when
		resp, err := http.Get(fmt.Sprintf("%s/accounts/%d", ts.Server.URL, createResponse.AccountID))
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.GetAccountResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, createResponse.AccountID, response.AccountID)
		assert.Equal(t, "11111111111", response.DocumentNumber)
	})

	t.Run("returns 404 when account does not exist", func(t *testing.T) {
		// when
		resp, err := http.Get(ts.Server.URL + "/accounts/999999")
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("returns 400 when account id is invalid", func(t *testing.T) {
		// when
		resp, err := http.Get(ts.Server.URL + "/accounts/invalid")
		require.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
