package domain

import (
	"errors"
	"time"
)

type Transaction struct {
	ID            int64
	AccountID     int64
	OperationType OperationType
	Amount        int64
	Currency      string
	Scale         int16
	EventDate     time.Time
}

var (
	ErrTransctionInvalidAccountID          = errors.New("invalid transaction account ID")
	ErrTranactionInvalidOperationTypeID    = errors.New("invalid transaction operation type ID")
	ErrTransactionInvalidTransactionAmount = errors.New("invalid transaction amount")
)

const (
	BRL      = "BRL"
	BRLScale = 2
)

func NewTransaction(
	accountID int64,
	operationType OperationType,
	amount int64,
) (*Transaction, []error) {
	errs := []error{}
	if accountID <= 0 {
		errs = append(errs, ErrTransctionInvalidAccountID)
	}
	if amount == 0 {
		errs = append(errs, ErrTransactionInvalidTransactionAmount)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	signedAmount := operationType.ApplySign(amount)

	now := time.Now()

	return &Transaction{
		ID:            0,
		AccountID:     accountID,
		OperationType: operationType,
		Amount:        signedAmount,
		Currency:      BRL,
		Scale:         BRLScale,
		EventDate:     now,
	}, nil
}

func (t *Transaction) IsDebit() bool {
	return t.Amount < 0
}

func (t *Transaction) IsCredit() bool {
	return t.Amount > 0
}
