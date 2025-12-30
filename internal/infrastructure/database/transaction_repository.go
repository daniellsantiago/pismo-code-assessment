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
		INSERT INTO transactions (account_id, operation_type_id, amount, event_date, balance) 
		VALUES ($1, $2, $3, $4, $5) 
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
		transaction.Balance,
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
		Balance:         transaction.Balance,
	}, nil
}

func (r *TransactionRepository) ListByAccountID(ctx context.Context, accountID int64) ([]*domain.Transaction, error) {
	query := `
		SELECT transaction_id, account_id, operation_type_id, amount, event_date, balance
		FROM transactions 
		WHERE account_id = $1
		ORDER BY event_date ASC
	`

	rows, err := r.db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		transaction := &domain.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.AccountID, &transaction.OperationTypeID, &transaction.Amount, &transaction.EventDate, &transaction.Balance)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *TransactionRepository) UpdateBalance(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	query := `
		UPDATE transactions 
		SET balance = $1 
		WHERE transaction_id = $2
		RETURNING transaction_id
	`

	err := r.db.QueryRowContext(ctx, query, transaction.Balance, transaction.ID).Scan(&transaction.ID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
