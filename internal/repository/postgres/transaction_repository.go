package postgres

import (
	"context"
	"database/sql"
	"mini-ledger/internal/domain"
)

func NewTransactionRepository(db *sql.DB) domain.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

type transactionRepository struct {
	db *sql.DB
}

// Create implements [domain.TransactionRepository].
func (t *transactionRepository) Create(
	ctx context.Context,
	transaction *domain.Transaction,
) (*domain.Transaction, error) {
	panic("unimplemented")
}
