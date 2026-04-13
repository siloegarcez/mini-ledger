package postgres

import (
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/internal/domain"
)

func mapOperationTypeDomainToModel(operationType *domain.OperationType) *model.OperationsTypes {
	return &model.OperationsTypes{
		OperationTypeID: operationType.OperationTypeID,
		Description:     operationType.Description,
		SignMultiplier:  operationType.SignMultiplier,
		CreatedAt:       operationType.CreatedAt,
	}
}

func mapOperationTypeModelToDomain(operationType *model.OperationsTypes) *domain.OperationType {
	return &domain.OperationType{
		OperationTypeID: operationType.OperationTypeID,
		Description:     operationType.Description,
		SignMultiplier:  operationType.SignMultiplier,
		CreatedAt:       operationType.CreatedAt,
	}
}
