package service

import (
	"context"
	"errors"
	"mini-ledger/internal/domain"
	"mini-ledger/internal/shared/dberrs"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-jet/jet/v2/qrm"
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

var (
	ErrInvalidAccountID     = huma.Error400BadRequest("invalid account id")
	ErrInvalidOperationType = huma.Error400BadRequest("operation type not found")
)

// Create implements [TransactionService].
func (t *transactionService) Create(
	ctx context.Context,
	cmd *CreateTransactionCommand,
) (*domain.Transaction, error) {
	transacType, err := t.operationTypeRepository.GetByID(ctx, cmd.OperationTypeID)
	if errors.Is(err, qrm.ErrNoRows) {
		return nil, ErrInvalidOperationType
	}
	if err != nil {
		return nil, err
	}
	amount, err := domain.NewBRLMoneyFromString(cmd.Amount)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}
	newTransac, err := domain.NewTransaction(
		cmd.AccountID,
		transacType.OperationTypeID,
		amount,
	)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	createdTransac, err := t.transactionRepository.Create(ctx, newTransac)
	if dberrs.Is(err, dberrs.ErrForeignKeyViolation) {
		return nil, ErrInvalidAccountID
	}
	if err != nil {
		return nil, err
	}

	return createdTransac, nil
}
