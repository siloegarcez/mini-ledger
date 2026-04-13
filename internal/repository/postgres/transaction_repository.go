package postgres

import (
	"context"
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/gen/dev/public/table"
	"mini-ledger/internal/domain"

	"github.com/rs/zerolog/log"
)

func NewTransactionRepository(executor Executor) domain.TransactionRepository {
	return &transactionRepository{
		executor: executor,
	}
}

type transactionRepository struct {
	executor Executor
}

// Create implements [domain.TransactionRepository].
func (t *transactionRepository) Create(
	ctx context.Context,
	transaction *domain.Transaction,
) (*domain.Transaction, error) {
	stmt := table.Transactions.INSERT(table.Transactions.MutableColumns.Except(table.Transactions.DefaultColumns)).
		MODEL(mapTransactionDomainToModel(transaction)).
		RETURNING(table.Transactions.AllColumns)

	log.Ctx(ctx).Debug().Msg(stmt.DebugSql())

	var dest model.Transactions

	err := stmt.QueryContext(ctx, t.executor, &dest)

	if err != nil {
		return nil, err
	}

	return mapTransactionModelToDomain(&dest)
}
