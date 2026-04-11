package domain

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// BRL -> 2,  2.50
// USD -> 2,  1.25
// JPY -> 0,  100 .
type Transaction struct {
	ID            int64
	AccountID     int64
	OperationType *OperationType
	Amount        int64
	Currency      string
	Scale         int16
	EventDate     time.Time
}

var (
	ErrInvalidTransctionAccountID       = errors.New("invalid transaction account ID")
	ErrInvalidTranactionOperationTypeID = errors.New("invalid transaction operation type ID")
	ErrInvalidTransactionAmount         = errors.New("invalid transaction amount")
	ErrInvalidTransactionCurrency       = errors.New("invalid transaction currency")
	ErrInvalidTransactionScale          = errors.New("invalid transaction scale")
)

const (
	TransactionCurrencyLength = 3
)

func NewTransaction(
	accountID int64,
	operationType *OperationType,
	amount int64,
	currency string,
	scale int16,
) (*Transaction, error) {
	if accountID <= 0 {
		return nil, ErrInvalidTransctionAccountID
	}
	if operationType == nil || !slices.Contains(AllOperationTypes(), operationType) {
		return nil, ErrInvalidTranactionOperationTypeID
	}
	if amount <= 0 {
		return nil, ErrInvalidTransactionAmount
	}
	if len(currency) != TransactionCurrencyLength {
		return nil, ErrInvalidTransactionCurrency
	}
	if scale < 0 {
		return nil, ErrInvalidTransactionScale
	}

	signedAmount := operationType.ApplySign(amount)

	now := time.Now()

	return &Transaction{
		ID:            0,
		AccountID:     accountID,
		OperationType: operationType,
		Amount:        signedAmount,
		Currency:      currency,
		Scale:         scale,
		EventDate:     now,
	}, nil
}

func (t *Transaction) IsDebit() bool {
	return t.Amount < 0
}

func (t *Transaction) IsCredit() bool {
	return t.Amount > 0
}

func (t *Transaction) Validate() []error {
	errs := []error{}
	if t.AccountID <= 0 {
		errs = append(errs, ErrInvalidTransctionAccountID)
	}
	if !slices.Contains(AllOperationTypes(), t.OperationType) {
		errs = append(errs, ErrInvalidTranactionOperationTypeID)
	}
	if t.Amount <= 0 {
		errs = append(errs, ErrInvalidTransactionAmount)
	}
	if len(t.Currency) != TransactionCurrencyLength {
		errs = append(errs, ErrInvalidTransactionCurrency)
	}
	if t.Scale < 0 {
		errs = append(errs, ErrInvalidTransactionScale)
	}

	if t.OperationType.IsCredit() && t.Amount < 0 {
		errs = append(
			errs,
			fmt.Errorf(
				"%w: debit transaction must have positive amount",
				ErrInvalidTransactionAmount,
			),
		)
	}

	if t.OperationType.IsDebit() && t.Amount > 0 {
		errs = append(
			errs,
			fmt.Errorf(
				"%w: credit transaction must have negative amount",
				ErrInvalidTransactionAmount,
			),
		)
	}

	return errs
}
