package postgres

import (
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/internal/domain"
)

func mapTransactionDomainToModel(transaction *domain.Transaction) *model.Transactions {
	return &model.Transactions{
		ID:              transaction.ID,
		AccountID:       transaction.AccountID,
		Currency:        transaction.Currency,
		Scale:           transaction.Amount.Scale(),
		Amount:          transaction.Amount.Int64(),
		OperationTypeID: transaction.OperationTypeID,
		EventDate:       transaction.EventDate,
	}
}

func mapTransactionModelToDomain(
	transaction *model.Transactions,
) (*domain.Transaction, error) {
	amount, err := domain.NewMoneyFromInt64(transaction.Amount, transaction.Scale)

	if err != nil {
		return nil, err
	}

	return &domain.Transaction{
		ID:              transaction.ID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Currency:        transaction.Currency,
		Amount:          amount,
		EventDate:       transaction.EventDate,
	}, nil
}
