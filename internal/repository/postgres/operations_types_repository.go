package postgres

import (
	"context"
	"database/sql"
	"mini-ledger/internal/domain"
)

func NewOperationTypeRepository(db *sql.DB) domain.OperationTypeRepository {
	return &operationTypeRepository{
		db: db,
	}
}

type operationTypeRepository struct {
	db *sql.DB
}

// GetByID implements [domain.OperationTypeRepository].
func (o *operationTypeRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.OperationType, error) {
	panic("unimplemented")
}
