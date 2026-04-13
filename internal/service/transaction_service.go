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
	transactor domain.Transactor,
) TransactionService {
	return &transactionService{
		transactionRepository:   transactionRepository,
		operationTypeRepository: operationTypeRepository,
		transactor:              transactor,
	}
}

type transactionService struct {
	transactionRepository   domain.TransactionRepository
	transactor              domain.Transactor
	operationTypeRepository domain.OperationTypeRepository
}

// Create implements [TransactionService].
func (t *transactionService) Create(
	ctx context.Context,
	cmd *CreateTransactionCommand,
) (*domain.Transaction, error) {
	panic("unimplemented")
}
