package domain

import "errors"

var (
	ErrInvalidDocumentNumber = errors.New("document number is required")
)

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
