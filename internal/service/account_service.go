package service

import (
	"context"
	"mini-ledger/internal/domain"
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

// Create implements [AccountService].
func (a *accountService) Create(
	ctx context.Context,
	documentNumber string,
) (*domain.Account, error) {
	panic("unimplemented")
}

// GetByID implements [AccountService].
func (a *accountService) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	panic("unimplemented")
}
