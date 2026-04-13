package domain

import (
	"errors"
	"time"
)

type Transaction struct {
	ID              int64
	AccountID       int64
	OperationTypeID int64
	Amount          BRLMoney
	Currency        string
	EventDate       time.Time
}

var (
	ErrTransactionInvalidAccountID         = errors.New("invalid transaction account ID")
	ErrTransactionInvalidOperationTypeID   = errors.New("invalid transaction operation type ID")
	ErrTransactionInvalidTransactionAmount = errors.New("invalid transaction amount")
	ErrTransactionNegativeAmount           = errors.New("invalid transaction amount sign")
)

const (
	BRL      = "BRL"
	BRLScale = 2
)

func NewTransaction(
	accountID int64,
	operationTypeID int64,
	amount BRLMoney,
) (*Transaction, error) {
	if accountID <= 0 {
		return nil, ErrTransactionInvalidAccountID
	}
	if amount.Int64() == 0 {
		return nil, ErrTransactionInvalidTransactionAmount
	}

	if operationTypeID <= 0 {
		return nil, ErrTransactionInvalidOperationTypeID
	}

	if amount.IsNegative() {
		return nil, ErrTransactionNegativeAmount
	}

	return &Transaction{
		ID:              0,
		AccountID:       accountID,
		OperationTypeID: operationTypeID,
		Amount:          amount,
		Currency:        BRL,
		EventDate:       time.Time{},
	}, nil
}

func (t *Transaction) IsDebit() bool {
	return t.Amount.IsNegative()
}

func (t *Transaction) IsCredit() bool {
	return t.Amount.IsPositive()
}
