package domain

import "time"

type Account struct {
	ID             int64
	DocumentNumber DocumentNumber
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewAccount(documentNumber string) (*Account, error) {
	docNum, err := NewDocumentNumber(documentNumber)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Account{
		ID:             0,
		DocumentNumber: docNum,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (a *Account) Validate() []error {
	errs := []error{}
	if a.DocumentNumber.IsEmpty() {
		errs = append(errs, ErrInvalidDocumentNumber)
	}
	return errs
}
