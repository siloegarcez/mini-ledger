package postgres

import (
	"context"
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/gen/dev/public/table"
	"mini-ledger/internal/domain"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/rs/zerolog/log"
)

func NewAccountsRepository(executor Executor) domain.AccountRepository {
	return &accountsRepository{
		executor: executor,
	}
}

type accountsRepository struct {
	executor Executor
}

// Create implements [domain.AccountRepository].
func (a *accountsRepository) Create(
	ctx context.Context,
	account *domain.Account,
) (*domain.Account, error) {
	stmt := table.Accounts.INSERT(table.Accounts.MutableColumns.Except(table.Accounts.DefaultColumns)).
		MODEL(mapDomainAccountToModel(account)).
		RETURNING(table.Accounts.AllColumns)

	log.Ctx(ctx).Debug().Msg(stmt.DebugSql())

	var dest model.Accounts

	err := stmt.QueryContext(ctx, a.executor, &dest)

	if err != nil {
		return nil, err
	}

	return mapAccountsModelToDomain(&dest)
}

// GetByID implements [domain.AccountRepository].
func (a *accountsRepository) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	stmt := table.Accounts.SELECT(table.Accounts.AllColumns).WHERE(table.Accounts.ID.EQ(Int64(id)))

	log.Ctx(ctx).Debug().Msg(stmt.DebugSql())

	var dest model.Accounts

	err := stmt.QueryContext(ctx, a.executor, &dest)

	if err != nil {
		return nil, err
	}

	return mapAccountsModelToDomain(&dest)
}
