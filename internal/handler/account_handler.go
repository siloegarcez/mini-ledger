package handler

import (
	"context"
	"mini-ledger/internal/service"
)

type AccountHandler interface {
	HandleCreateAccount(
		ctx context.Context,
		req *accountCreateRequest,
	) (*accountCreateResponse, error)
	HandleGetAccountByID(
		ctx context.Context,
		req *accountGetByIDRequest,
	) (*accountGetByIDResponse, error)
}

func NewAccountHandler(service service.AccountService) AccountHandler {
	return &accountHandler{service: service}
}

type accountHandler struct {
	service service.AccountService
}

// HandleCreateAccount implements [AccountHandler].
func (a *accountHandler) HandleCreateAccount(
	ctx context.Context,
	req *accountCreateRequest,
) (*accountCreateResponse, error) {
	panic("unimplemented")
}

// HandleGetAccountByID implements [AccountHandler].
func (a *accountHandler) HandleGetAccountByID(
	ctx context.Context,
	req *accountGetByIDRequest,
) (*accountGetByIDResponse, error) {
	panic("unimplemented")
}
