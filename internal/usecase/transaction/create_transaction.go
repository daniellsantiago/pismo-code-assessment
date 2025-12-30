package transaction

import (
	"context"

	"github.com/nubank/pismo-code-assessment/internal/domain"
)

type CreateTransaction struct {
	repo domain.TransactionRepository
}

func NewCreateTransaction(repo domain.TransactionRepository) *CreateTransaction {
	return &CreateTransaction{repo: repo}
}

func (c *CreateTransaction) Execute(ctx context.Context, accountID int64, operationTypeID int, amount float64) (*domain.Transaction, error) {
	var balance float64
	operationType := domain.OperationType(operationTypeID)
	if operationType.IsDebit() {
		balance = -amount
	} else {
		balance = amount

		pastTransactions, err := c.repo.ListByAccountID(ctx, accountID)
		if err != nil {
			return nil, err
		}

		for _, pastTransaction := range pastTransactions {
			if pastTransaction.OperationTypeID.IsDebit() && pastTransaction.IsNegative() && balance > 0 {
				if balance >= -pastTransaction.Balance {
					balance = balance + pastTransaction.Balance
					pastTransaction.Balance = 0
				} else {
					pastTransaction.Balance = pastTransaction.Balance + balance
					balance = 0
				}

				_, err := c.repo.UpdateBalance(ctx, pastTransaction)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	transaction, err := domain.NewTransaction(accountID, domain.OperationType(operationTypeID), amount, balance)
	if err != nil {
		return nil, err
	}

	return c.repo.Create(ctx, transaction)
}
