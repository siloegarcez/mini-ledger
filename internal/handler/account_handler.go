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

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{accountService: accountService}
}

type accountHandler struct {
	accountService service.AccountService
}

// HandleCreateAccount implements [AccountHandler].
func (a *accountHandler) HandleCreateAccount(
	ctx context.Context,
	req *accountCreateRequest,
) (*accountCreateResponse, error) {
	acc, err := a.accountService.Create(ctx, req.Body.DocumentNumber)

	if err != nil {
		return nil, err
	}

	return &accountCreateResponse{
		Body: mapDomainAccountToAccountCreateResponse(acc),
	}, nil
}

// HandleGetAccountByID implements [AccountHandler].
func (a *accountHandler) HandleGetAccountByID(
	ctx context.Context,
	req *accountGetByIDRequest,
) (*accountGetByIDResponse, error) {
	acc, err := a.accountService.GetByID(ctx, req.AccountID)

	if err != nil {
		return nil, err
	}

	return &accountGetByIDResponse{
		Body: mapDomainAccountToAccountGetByIDResponse(acc),
	}, nil
}
