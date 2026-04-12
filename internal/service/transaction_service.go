package service

import (
	"context"
	"mini-ledger/internal/domain"
)

type TransactionService interface {
	Create(ctx context.Context, cmd *CreateTransactionCommand) (*domain.Transaction, error)
}

func NewTransactionService(
	transactionRepository domain.TransactionRepository,
	operationTypeRepository domain.OperationTypeRepository,
) TransactionService {
	return &transactionService{
		transactionRepository:   transactionRepository,
		operationTypeRepository: operationTypeRepository,
	}
}

type transactionService struct {
	transactionRepository   domain.TransactionRepository
	operationTypeRepository domain.OperationTypeRepository
}

// Create implements [TransactionService].
func (t *transactionService) Create(
	ctx context.Context,
	cmd *CreateTransactionCommand,
) (*domain.Transaction, error) {
	panic("unimplemented")
}
