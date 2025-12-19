package domain

import (
	"context"
	"time"
)

type OperationType int

const (
	OperationTypePurchase            OperationType = 1
	OperationTypeInstallmentPurchase OperationType = 2
	OperationTypeWithdrawal          OperationType = 3
	OperationTypePayment             OperationType = 4
)

func (o OperationType) IsValid() bool {
	switch o {
	case OperationTypePurchase, OperationTypeInstallmentPurchase, OperationTypeWithdrawal, OperationTypePayment:
		return true
	}
	return false
}

func (o OperationType) IsDebit() bool {
	return o == OperationTypePurchase || o == OperationTypeInstallmentPurchase || o == OperationTypeWithdrawal
}

//go:generate mockgen -source=transaction.go -destination=mocks/transaction_mock.go -package=mocks
type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) (*Transaction, error)
}

type Transaction struct {
	ID              int64
	AccountID       int64
	OperationTypeID OperationType
	Amount          float64
	EventDate       time.Time
}

func NewTransaction(accountID int64, operationTypeID OperationType, amount float64) (*Transaction, error) {
	if !operationTypeID.IsValid() {
		return nil, ErrInvalidOperationType
	}

	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if operationTypeID.IsDebit() {
		amount = -amount
	}

	return &Transaction{
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
		EventDate:       time.Now(),
	}, nil
}
