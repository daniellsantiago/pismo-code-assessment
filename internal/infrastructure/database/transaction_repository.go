package database

import (
	"context"
	"database/sql"
	"strings"

	"github.com/lib/pq"
	"github.com/nubank/pismo-code-assessment/internal/domain"
)

const foreignKeyViolationCode = "23503"

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	query := `
		INSERT INTO transactions (account_id, operation_type_id, amount, event_date) 
		VALUES ($1, $2, $3, $4) 
		RETURNING transaction_id
	`

	var id int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		transaction.AccountID,
		transaction.OperationTypeID,
		transaction.Amount,
		transaction.EventDate,
	).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == foreignKeyViolationCode {
			if strings.Contains(pgErr.Constraint, "account") {
				return nil, domain.ErrAccountNotFound
			}
			if strings.Contains(pgErr.Constraint, "operation_type") {
				return nil, domain.ErrInvalidOperationType
			}
		}
		return nil, err
	}

	return &domain.Transaction{
		ID:              id,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate,
	}, nil
}
