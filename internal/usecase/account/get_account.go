package account

import (
	"context"

	"github.com/nubank/pismo-code-assessment/internal/domain"
)

type GetAccount struct {
	repo domain.AccountRepository
}

func NewGetAccount(repo domain.AccountRepository) *GetAccount {
	return &GetAccount{repo: repo}
}

func (g *GetAccount) Execute(ctx context.Context, accountID int64) (*domain.Account, error) {
	return g.repo.FindByID(ctx, accountID)
}
