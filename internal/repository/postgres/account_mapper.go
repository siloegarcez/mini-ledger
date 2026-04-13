package postgres

import (
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/internal/domain"
)

func mapDomainAccountToModel(account *domain.Account) *model.Accounts {
	return &model.Accounts{
		ID:             account.ID,
		DocumentNumber: account.DocumentNumber.String(),
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}
}

func mapAccountsModelToDomain(account *model.Accounts) (*domain.Account, error) {
	documentNumber, err := domain.NewDocumentNumber(account.DocumentNumber)
	if err != nil {
		return nil, err
	}
	return &domain.Account{
		ID:             account.ID,
		DocumentNumber: documentNumber,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}, nil
}
