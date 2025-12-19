package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperationType_IsValid(t *testing.T) {
	t.Run("returns true for valid operation types", func(t *testing.T) {
		validTypes := []OperationType{
			OperationTypePurchase,
			OperationTypeInstallmentPurchase,
			OperationTypeWithdrawal,
			OperationTypePayment,
		}

		for _, opType := range validTypes {
			assert.True(t, opType.IsValid())
		}
	})

	t.Run("returns false for invalid operation types", func(t *testing.T) {
		invalidTypes := []OperationType{0, 5, 100, -1}

		for _, opType := range invalidTypes {
			assert.False(t, opType.IsValid())
		}
	})
}

func TestOperationType_IsDebit(t *testing.T) {
	t.Run("returns true for debit operations", func(t *testing.T) {
		debitTypes := []OperationType{
			OperationTypePurchase,
			OperationTypeInstallmentPurchase,
			OperationTypeWithdrawal,
		}

		for _, opType := range debitTypes {
			assert.True(t, opType.IsDebit())
		}
	})

	t.Run("returns false for credit operations", func(t *testing.T) {
		assert.False(t, OperationTypePayment.IsDebit())
	})
}

func TestNewTransaction(t *testing.T) {
	t.Run("creates transaction with negative amount for purchase", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationTypePurchase, 50.0)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), transaction.AccountID)
		assert.Equal(t, OperationTypePurchase, transaction.OperationTypeID)
		assert.Equal(t, -50.0, transaction.Amount)
	})

	t.Run("creates transaction with negative amount for installment purchase", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationTypeInstallmentPurchase, 100.0)

		assert.NoError(t, err)
		assert.Equal(t, -100.0, transaction.Amount)
	})

	t.Run("creates transaction with negative amount for withdrawal", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationTypeWithdrawal, 25.0)

		assert.NoError(t, err)
		assert.Equal(t, -25.0, transaction.Amount)
	})

	t.Run("creates transaction with positive amount for payment", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationTypePayment, 123.45)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), transaction.AccountID)
		assert.Equal(t, OperationTypePayment, transaction.OperationTypeID)
		assert.Equal(t, 123.45, transaction.Amount)
	})

	t.Run("returns error for invalid operation type", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationType(99), 50.0)

		assert.Nil(t, transaction)
		assert.ErrorIs(t, err, ErrInvalidOperationType)
	})

	t.Run("returns error when amount is zero", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationTypePurchase, 0)

		assert.Nil(t, transaction)
		assert.ErrorIs(t, err, ErrInvalidAmount)
	})

	t.Run("returns error when amount is negative", func(t *testing.T) {
		transaction, err := NewTransaction(1, OperationTypePurchase, -50.0)

		assert.Nil(t, transaction)
		assert.ErrorIs(t, err, ErrInvalidAmount)
	})
}
