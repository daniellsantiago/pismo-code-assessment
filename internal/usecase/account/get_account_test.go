package account

import (
	"context"
	"errors"
	"testing"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRepo := mocks.NewMockAccountRepository(ctrl)
	usecase := NewGetAccount(mockedRepo)

	t.Run("retrieves account successfully", func(t *testing.T) {
		// given
		accountID := int64(1)
		documentNumber := "12345678900"
		expectedAccount := &domain.Account{
			ID:             accountID,
			DocumentNumber: documentNumber,
		}

		// when
		mockedRepo.EXPECT().FindByID(gomock.Any(), accountID).Return(expectedAccount, nil)

		account, err := usecase.Execute(context.Background(), accountID)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, accountID, account.ID)
		assert.Equal(t, documentNumber, account.DocumentNumber)
	})

	t.Run("returns error when repository fails to fetch account", func(t *testing.T) {
		// given
		accountID := int64(999)

		// when
		mockedRepo.EXPECT().FindByID(gomock.Any(), accountID).Return(nil, errors.New("repository error"))

		account, err := usecase.Execute(context.Background(), accountID)

		// then
		assert.Nil(t, account)
		assert.Error(t, err)
	})
}
