package domain

import (
	"errors"
	"time"
)

type Transaction struct {
	ID            int64
	AccountID     int64
	OperationType OperationType
	Amount        BRLMoney
	Currency      string
	EventDate     time.Time
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
	operationType OperationType,
	amount BRLMoney,
) (*Transaction, []error) {
	errs := []error{}
	if accountID <= 0 {
		errs = append(errs, ErrTransactionInvalidAccountID)
	}
	if amount.Int64() == 0 {
		errs = append(errs, ErrTransactionInvalidTransactionAmount)
	}

	if operationType.OperationTypeID <= 0 {
		errs = append(errs, ErrTransactionInvalidOperationTypeID)
	}

	if amount.IsNegative() {
		errs = append(errs, ErrTransactionNegativeAmount)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	if operationType.IsDebit() {
		amount = amount.Neg()
	}

	now := time.Now()

	return &Transaction{
		ID:            0,
		AccountID:     accountID,
		OperationType: operationType,
		Amount:        amount,
		Currency:      BRL,
		EventDate:     now,
	}, nil
}

func (t *Transaction) IsDebit() bool {
	return t.Amount.IsNegative()
}

func (t *Transaction) IsCredit() bool {
	return t.Amount.IsPositive()
}
