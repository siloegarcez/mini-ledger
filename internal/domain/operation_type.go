package domain

import (
	"errors"
	"fmt"
	"slices"
)

type OperationType struct {
	OperationID    int16
	Description    string
	SignMultiplier int16
}

var (
	OperationTypePurchase = &OperationType{
		OperationID:    1,
		Description:    "PURCHASE",
		SignMultiplier: -1,
	}
	OperationTypeInstallmentPurchase = &OperationType{
		OperationID:    2,
		Description:    "INSTALLMENT_PURCHASE",
		SignMultiplier: -1,
	}
	OperationTypeWithdrawal = &OperationType{
		OperationID:    3,
		Description:    "WITHDRAWAL",
		SignMultiplier: -1,
	}
	OperationTypePayment = &OperationType{
		OperationID:    4,
		Description:    "PAYMENT",
		SignMultiplier: 1,
	}
)

var (
	ErrOperationTypeNotFound = errors.New("operation type not found")
	ErrInvalidOperationType  = errors.New("invalid operation type")
)

func AllOperationTypes() []*OperationType {
	return []*OperationType{
		OperationTypePurchase,
		OperationTypeInstallmentPurchase,
		OperationTypeWithdrawal,
		OperationTypePayment,
	}
}

func GetOperationTypeByID(id int16) (*OperationType, error) {
	for _, ot := range AllOperationTypes() {
		if ot.OperationID == id {
			return ot, nil
		}
	}
	return nil, ErrOperationTypeNotFound
}

func GetDebitOperationTypes() []*OperationType {
	debitTypes := []*OperationType{}
	for _, ot := range AllOperationTypes() {
		if ot.IsDebit() {
			debitTypes = append(debitTypes, ot)
		}
	}
	return debitTypes
}

func GetCreditOperationTypes() []*OperationType {
	debitTypes := []*OperationType{}
	for _, ot := range AllOperationTypes() {
		if ot.IsCredit() {
			debitTypes = append(debitTypes, ot)
		}
	}
	return debitTypes
}

func (ot *OperationType) ApplySign(amount int64) int64 {
	if amount < 0 {
		amount = -amount
	}
	return amount * int64(ot.SignMultiplier)
}

func (ot *OperationType) IsDebit() bool {
	return ot.SignMultiplier == -1
}

func (ot *OperationType) IsCredit() bool {
	return ot.SignMultiplier == 1
}

func (ot *OperationType) Validate() []error {
	errs := []error{}
	if ot.OperationID <= 0 {
		errs = append(
			errs,
			fmt.Errorf("%w: operation id must be positive", ErrInvalidOperationType),
		)
	}

	if !slices.Contains(AllOperationTypes(), ot) {
		errs = append(
			errs,
			fmt.Errorf(
				"%w: operation type with id %d does not exist",
				ErrInvalidOperationType,
				ot.OperationID,
			),
		)
	}

	if ot.Description == "" {
		errs = append(errs, fmt.Errorf("%w: description cannot be empty", ErrInvalidOperationType))
	}

	if ot.SignMultiplier != -1 && ot.SignMultiplier != 1 {
		errs = append(
			errs,
			fmt.Errorf("%w: sign multiplier must be either -1 or 1", ErrInvalidOperationType),
		)
	}

	return errs
}
