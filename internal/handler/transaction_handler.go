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
	transac, err := t.transactionService.Create(ctx, &service.CreateTransactionCommand{
		AccountID:       req.Body.AccountID,
		OperationTypeID: req.Body.OperationTypeID,
		Amount:          string(req.Body.Amount),
	})

	if err != nil {
		return nil, err
	}

	return &transactionCreateResponse{
		Body: mapDomainTransactionToTransactionCreateResponse(transac),
	}, nil
}
