package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	t.Run("creates account when provided document number is not empty", func(t *testing.T) {
		account, err := NewAccount("12345678900")

		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, "12345678900", account.DocumentNumber)
		assert.Equal(t, int64(0), account.ID)
	})

	t.Run("returns error when document number is empty", func(t *testing.T) {
		account, err := NewAccount("")

		assert.Nil(t, account)
		assert.ErrorIs(t, err, ErrInvalidDocumentNumber)
	})
}
