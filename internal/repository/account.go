package repository

import (
	"context"

	"github.com/nubank/pismo-code-assessment/internal/domain"
)

//go:generate mockgen -source=account.go -destination=mocks/account_mock.go -package=mocks
type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) (*domain.Account, error)
}
