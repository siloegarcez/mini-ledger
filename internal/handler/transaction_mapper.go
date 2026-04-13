package handler

import (
	"mini-ledger/internal/domain"
)

func mapDomainTransactionToTransactionCreateResponse(
	transaction *domain.Transaction,
) TransactionCreateResponse {
	return TransactionCreateResponse{
		TransactionID:   transaction.ID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          Number(transaction.Amount.String()),
	}
}
