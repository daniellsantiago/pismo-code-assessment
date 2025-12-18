package account

import (
	"context"

	"github.com/nubank/pismo-code-assessment/internal/domain"
)

type CreateAccount struct {
	repo domain.AccountRepository
}

func NewCreateAccount(repo domain.AccountRepository) *CreateAccount {
	return &CreateAccount{repo: repo}
}

func (c *CreateAccount) Execute(ctx context.Context, documentNumber string) (*domain.Account, error) {
	account, err := domain.NewAccount(documentNumber)
	if err != nil {
		return nil, err
	}

	return c.repo.Create(ctx, account)
}
