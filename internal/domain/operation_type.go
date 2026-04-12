package domain

import (
	"errors"
	"strings"
)

type OperationType struct {
	OperationTypeID int64
	Description     string
	SignMultiplier  int16
}

const (
	MaxDescriptionLength = 30
	CreditSignMultiplier = 1
	DebitSignMultiplier  = -1
)

var (
	ErrOperationTypeInvalidSignMultiplier  = errors.New("invalid operation type sign multiplier")
	ErrOperationTypeInvalidOperationTypeID = errors.New("invalid operation type ID")
	ErrOperationTypeEmptyDescription       = errors.New("empty operation type description")
	ErrOperationTypeInvalidDescriptionLen  = errors.New("invalid operation type description length")
)

func NewOperationType(
	operationTypeID int64,
	description string,
	signMultiplier int16,
) (*OperationType, []error) {
	errs := []error{}

	if operationTypeID <= 0 {
		errs = append(errs, ErrOperationTypeInvalidOperationTypeID)
	}
	if signMultiplier != CreditSignMultiplier && signMultiplier != DebitSignMultiplier {
		errs = append(errs, ErrOperationTypeInvalidSignMultiplier)
	}
	if len(strings.TrimSpace(description)) == 0 {
		errs = append(errs, ErrOperationTypeEmptyDescription)
	}
	if len(description) > MaxDescriptionLength {
		errs = append(errs, ErrOperationTypeInvalidDescriptionLen)
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &OperationType{
		OperationTypeID: operationTypeID,
		Description:     description,
		SignMultiplier:  signMultiplier,
	}, nil
}

func (ot *OperationType) ApplySign(amount int64) int64 {
	if amount < 0 {
		amount = -amount
	}
	return amount * int64(ot.SignMultiplier)
}

func (ot *OperationType) IsDebit() bool {
	return ot.SignMultiplier == DebitSignMultiplier
}

func (ot *OperationType) IsCredit() bool {
	return ot.SignMultiplier == CreditSignMultiplier
}
