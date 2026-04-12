package domain

import (
	"context"
)

type AccountRepository interface {
	Create(ctx context.Context, account *Account) (*Account, error)
	GetByID(ctx context.Context, id int64) (*Account, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) (*Transaction, error)
}

type OperationsTypesRepository interface {
	GetByID(ctx context.Context, id int64) (*OperationType, error)
}
