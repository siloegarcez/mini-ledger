package domain

import (
	"errors"
	"strings"
	"time"
)

type OperationType struct {
	OperationTypeID int64
	Description     string
	SignMultiplier  int16
	CreatedAt       time.Time
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
) (*OperationType, error) {
	description = strings.TrimSpace(description)
	if operationTypeID <= 0 {
		return nil, ErrOperationTypeInvalidOperationTypeID
	}
	if signMultiplier != CreditSignMultiplier && signMultiplier != DebitSignMultiplier {
		return nil, ErrOperationTypeInvalidSignMultiplier
	}
	if len(description) == 0 {
		return nil, ErrOperationTypeEmptyDescription
	}
	if len(description) > MaxDescriptionLength {
		return nil, ErrOperationTypeInvalidDescriptionLen
	}

	return &OperationType{
		OperationTypeID: operationTypeID,
		Description:     description,
		SignMultiplier:  signMultiplier,
		CreatedAt:       time.Time{},
	}, nil
}

func (ot *OperationType) IsDebit() bool {
	return ot.SignMultiplier == DebitSignMultiplier
}

func (ot *OperationType) IsCredit() bool {
	return ot.SignMultiplier == CreditSignMultiplier
}
