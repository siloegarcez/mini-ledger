package handler

import "mini-ledger/internal/domain"

func mapDomainAccountToAccountCreateResponse(acc *domain.Account) AccountCreateResponse {
	return AccountCreateResponse{
		AccountID:      acc.ID,
		DocumentNumber: acc.DocumentNumber.String(),
	}
}

func mapDomainAccountToAccountGetByIDResponse(acc *domain.Account) AccountGetByIDResponse {
	return AccountGetByIDResponse{
		AccountID:      acc.ID,
		DocumentNumber: acc.DocumentNumber.String(),
	}
}
