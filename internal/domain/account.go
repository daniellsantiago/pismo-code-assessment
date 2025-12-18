package domain

import (
	"context"
)

//go:generate mockgen -source=account.go -destination=mocks/account_mock.go -package=mocks
type AccountRepository interface {
	Create(ctx context.Context, account *Account) (*Account, error)
	FindByID(ctx context.Context, ID int64) (*Account, error)
}

type Account struct {
	ID             int64
	DocumentNumber string
}

func NewAccount(documentNumber string) (*Account, error) {
	if documentNumber == "" {
		return nil, ErrInvalidDocumentNumber
	}

	return &Account{
		DocumentNumber: documentNumber,
	}, nil
}
