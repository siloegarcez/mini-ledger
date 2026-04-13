package postgres

import (
	"context"
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/gen/dev/public/table"
	"mini-ledger/internal/domain"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/rs/zerolog/log"
)

func NewOperationTypeRepository(executor Executor) domain.OperationTypeRepository {
	return &operationTypeRepository{
		executor: executor,
	}
}

type operationTypeRepository struct {
	executor Executor
}

// GetByID implements [domain.OperationTypeRepository].
func (o *operationTypeRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.OperationType, error) {
	stmt := table.OperationsTypes.SELECT(table.OperationsTypes.AllColumns).
		WHERE(table.OperationsTypes.OperationTypeID.EQ(Int64(id)))

	log.Ctx(ctx).Debug().Msg(stmt.DebugSql())

	var dest model.OperationsTypes

	err := stmt.QueryContext(ctx, o.executor, &dest)

	if err != nil {
		return nil, err
	}

	return mapOperationTypeModelToDomain(&dest), nil
}
