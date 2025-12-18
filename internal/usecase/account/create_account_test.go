package account

import (
	"context"
	"testing"

	"github.com/nubank/pismo-code-assessment/internal/domain"
	"github.com/nubank/pismo-code-assessment/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedRepo := mocks.NewMockAccountRepository(ctrl)
	usecase := NewCreateAccount(mockedRepo)

	t.Run("creates account successfully", func(t *testing.T) {
		// given
		documentNumber := "12345678900"
		accoutID := int64(1)
		expectedAccount := &domain.Account{
			ID:             accoutID,
			DocumentNumber: documentNumber,
		}

		// when
		mockedRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(expectedAccount, nil)

		account, err := usecase.Execute(context.Background(), "12345678900")

		// then
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, accoutID, account.ID)
		assert.Equal(t, documentNumber, account.DocumentNumber)
	})

	t.Run("returns error when document number is empty", func(t *testing.T) {
		// given
		documentNumber := ""

		// when
		account, err := usecase.Execute(context.Background(), documentNumber)

		// then
		assert.Nil(t, account)
		assert.ErrorIs(t, err, domain.ErrInvalidDocumentNumber)
	})
}
