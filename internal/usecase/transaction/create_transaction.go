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
	transaction, err := domain.NewTransaction(accountID, domain.OperationType(operationTypeID), amount)
	if err != nil {
		return nil, err
	}

	return c.repo.Create(ctx, transaction)
}
