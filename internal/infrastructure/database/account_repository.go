package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/nubank/pismo-code-assessment/internal/domain"
)

const uniqueViolationCode = "23505"

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	query := `INSERT INTO accounts (document_number) VALUES ($1) RETURNING account_id`

	var id int64
	err := r.db.QueryRowContext(ctx, query, account.DocumentNumber).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == uniqueViolationCode {
			return nil, domain.ErrAccountAlreadyExists
		}
		return nil, err
	}

	return &domain.Account{
		ID:             id,
		DocumentNumber: account.DocumentNumber,
	}, nil
}

func (r *AccountRepository) FindByID(ctx context.Context, ID int64) (*domain.Account, error) {
	var account domain.Account

	query := `SELECT account_id, document_number FROM accounts WHERE account_id = ($1)`

	err := r.db.QueryRowContext(ctx, query, ID).Scan(
		&account.ID,
		&account.DocumentNumber,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrAccountNotFound
		}

		return nil, err
	}

	return &account, err
}
