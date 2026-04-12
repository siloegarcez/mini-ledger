package handler

import (
	"context"
	"mini-ledger/internal/service"
)

type TransactionHandler interface {
	HandleCreateTransaction(
		ctx context.Context,
		req *transactionCreateRequest,
	) (*transactionCreateResponse, error)
}

func NewTransactionHandler(service service.TransactionService) TransactionHandler {
	return &transactionHandler{transactionService: service}
}

type transactionHandler struct {
	transactionService service.TransactionService
}

// HandleCreateTransaction implements [TransactionHandler].
func (t *transactionHandler) HandleCreateTransaction(
	ctx context.Context,
	req *transactionCreateRequest,
) (*transactionCreateResponse, error) {
	panic("unimplemented")
}
