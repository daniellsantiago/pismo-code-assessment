package transaction

import (
	"context"
	"errors"
	"testing"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockTransactionRepository(ctrl)
	usecase := NewCreateTransaction(mockRepo)

	t.Run("creates transaction successfully", func(t *testing.T) {
		// given
		accountID := int64(1)
		operationTypeID := 4 // Payment
		amount := 123.45

		// when
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
				tx.ID = 1
				return tx, nil
			},
		)

		transaction, err := usecase.Execute(context.Background(), accountID, operationTypeID, amount)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, transaction)
		assert.Equal(t, int64(1), transaction.ID)
		assert.Equal(t, accountID, transaction.AccountID)
		assert.Equal(t, domain.OperationTypePayment, transaction.OperationTypeID)
		assert.Equal(t, 123.45, transaction.Amount)
	})

	t.Run("creates transaction with negative amount for purchase", func(t *testing.T) {
		// given
		accountID := int64(1)
		operationTypeID := 1 // Purchase
		amount := 50.0

		// when
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
				tx.ID = 2
				return tx, nil
			},
		)

		transaction, err := usecase.Execute(context.Background(), accountID, operationTypeID, amount)

		// then
		assert.NoError(t, err)
		assert.Equal(t, -50.0, transaction.Amount)
	})

	t.Run("returns error when account not found", func(t *testing.T) {
		// given
		accountID := int64(999)

		// when
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, domain.ErrAccountNotFound)

		transaction, err := usecase.Execute(context.Background(), accountID, 1, 50.0)

		// then
		assert.Nil(t, transaction)
		assert.ErrorIs(t, err, domain.ErrAccountNotFound)
	})

	t.Run("returns error when operation type is invalid", func(t *testing.T) {
		// given
		accountID := int64(1)

		// when - domain validation catches invalid operation type
		transaction, err := usecase.Execute(context.Background(), accountID, 99, 50.0)

		// then
		assert.Nil(t, transaction)
		assert.ErrorIs(t, err, domain.ErrInvalidOperationType)
	})

	t.Run("returns error when amount is invalid", func(t *testing.T) {
		// given
		accountID := int64(1)

		// when
		transaction, err := usecase.Execute(context.Background(), accountID, 1, 0)

		// then
		assert.Nil(t, transaction)
		assert.ErrorIs(t, err, domain.ErrInvalidAmount)
	})

	t.Run("returns error when repository fails", func(t *testing.T) {
		// given
		accountID := int64(1)

		// when
		mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("database error"))

		transaction, err := usecase.Execute(context.Background(), accountID, 1, 50.0)

		// then
		assert.Nil(t, transaction)
		assert.Error(t, err)
	})
}
