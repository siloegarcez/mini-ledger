package domain

import (
	"errors"
	"time"
)

type Account struct {
	ID             int64
	DocumentNumber DocumentNumber
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var (
	ErrAccountDocumentNumberEmpty = errors.New("empty document number")
)

func NewAccount(documentNumber DocumentNumber) (*Account, error) {
	if documentNumber.String() == "" {
		return nil, ErrAccountDocumentNumberEmpty
	}

	return &Account{
		ID:             0,
		DocumentNumber: documentNumber,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
	}, nil
}
