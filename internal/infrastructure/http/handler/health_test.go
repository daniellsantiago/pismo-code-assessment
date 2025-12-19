package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler_Check(t *testing.T) {
	t.Run("returns healthy when database is connected", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectPing()

		handler := NewHealthHandler(db)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.Check(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response HealthResponse
		err = json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "connected", response.Database)
	})

	t.Run("returns unhealthy when database is disconnected", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
		require.NoError(t, err)

		mock.ExpectPing().WillReturnError(assert.AnError)
		db.Close() // Close the connection to simulate disconnect

		handler := NewHealthHandler(db)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.Check(rec, req)

		assert.Equal(t, http.StatusServiceUnavailable, rec.Code)

		var response HealthResponse
		err = json.NewDecoder(rec.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, "unhealthy", response.Status)
		assert.Equal(t, "disconnected", response.Database)
	})
}
