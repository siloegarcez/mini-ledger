package service

import (
	"context"
	"errors"
	"mini-ledger/internal/domain"
	"mini-ledger/internal/shared/dberrs"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-jet/jet/v2/qrm"
)

type AccountService interface {
	Create(ctx context.Context, documentNumber string) (*domain.Account, error)
	GetByID(ctx context.Context, id int64) (*domain.Account, error)
}

func NewAccountService(
	accountRepository domain.AccountRepository,
	transactor domain.Transactor,
) AccountService {
	return &accountService{
		accountRepository: accountRepository,
		transactor:        transactor,
	}
}

type accountService struct {
	accountRepository domain.AccountRepository
	transactor        domain.Transactor
}

var (
	ErrDuplicateDocumentNumber = huma.Error409Conflict(
		"document number already exists",
	)
	ErrAccountNotFound = huma.Error404NotFound("account not found")
)

// Create implements [AccountService].
func (a *accountService) Create(
	ctx context.Context,
	documentNumber string,
) (*domain.Account, error) {
	docNum, err := domain.NewDocumentNumber(documentNumber)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	newAcc, err := domain.NewAccount(docNum)
	if err != nil {
		return nil, huma.Error422UnprocessableEntity(err.Error())
	}

	createdAcc, err := a.accountRepository.Create(ctx, newAcc)

	if dberrs.Is(err, dberrs.ErrDuplicateKeyViolation) {
		return nil, ErrDuplicateDocumentNumber
	}

	if err != nil {
		return nil, err
	}

	return createdAcc, nil
}

// GetByID implements [AccountService].
func (a *accountService) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	acc, err := a.accountRepository.GetByID(ctx, id)
	if errors.Is(err, qrm.ErrNoRows) {
		return nil, ErrAccountNotFound
	}
	if err != nil {
		return nil, err
	}
	return acc, nil
}
