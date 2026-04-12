package domain

import "context"

type AccountService interface {
	CreateAccount(ctx context.Context, documentNumber DocumentNumber) (*Account, error)
	GetAccountByID(ctx context.Context, id int64) (*Account, error)
}
