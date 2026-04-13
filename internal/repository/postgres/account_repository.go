package postgres

import (
	"context"
	"database/sql"
	"mini-ledger/internal/domain"
)

func NewAccountsRepository(db *sql.DB) domain.AccountRepository {
	return &accountsRepository{
		db: db,
	}
}

type accountsRepository struct {
	db *sql.DB
}

// Create implements [domain.AccountRepository].
func (a *accountsRepository) Create(
	ctx context.Context,
	account *domain.Account,
) (*domain.Account, error) {
	panic("unimplemented")
}

// GetByID implements [domain.AccountRepository].
func (a *accountsRepository) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	panic("unimplemented")
}
