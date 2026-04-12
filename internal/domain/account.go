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

func NewAccount(documentNumber DocumentNumber) (*Account, []error) {
	errs := []error{}

	if documentNumber.String() == "" {
		errs = append(errs, ErrAccountDocumentNumberEmpty)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	now := time.Now()

	return &Account{
		ID:             0,
		DocumentNumber: documentNumber,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
