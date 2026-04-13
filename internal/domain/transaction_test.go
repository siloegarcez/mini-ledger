package domain_test

import (
	"mini-ledger/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransaction(t *testing.T) {
	validCreditOperationType := domain.OperationType{
		OperationTypeID: 4,
		Description:     "PAYMENT",
		SignMultiplier:  domain.CreditSignMultiplier,
	}

	validDebitOperationType := domain.OperationType{
		OperationTypeID: 1,
		Description:     "PURCHASE",
		SignMultiplier:  domain.DebitSignMultiplier,
	}

	positiveAmount, err := domain.NewBRLMoneyFromString("150.25")
	require.NoError(t, err)

	zeroAmount, err := domain.NewBRLMoneyFromString("0")
	require.NoError(t, err)

	negativeAmount, err := domain.NewBRLMoneyFromString("-20.10")
	require.NoError(t, err)

	tests := []struct {
		name          string
		accountID     int64
		operationType domain.OperationType
		amount        domain.BRLMoney
		wantErrs      []error
	}{
		{
			name:          "valid credit transaction keeps positive amount",
			accountID:     10,
			operationType: validCreditOperationType,
			amount:        positiveAmount,
			wantErrs:      nil,
		},
		{
			name:          "valid debit transaction negates amount",
			accountID:     20,
			operationType: validDebitOperationType,
			amount:        positiveAmount,
			wantErrs:      nil,
		},
		{
			name:          "invalid account id returns error",
			accountID:     0,
			operationType: validCreditOperationType,
			amount:        positiveAmount,
			wantErrs:      []error{domain.ErrTransactionInvalidAccountID},
		},
		{
			name:      "invalid operation type id returns error",
			accountID: 1,
			operationType: domain.OperationType{ //nolint: exhaustruct
				SignMultiplier: domain.CreditSignMultiplier,
			},
			amount:   positiveAmount,
			wantErrs: []error{domain.ErrTransactionInvalidOperationTypeID},
		},
		{
			name:          "zero amount returns error",
			accountID:     1,
			operationType: validCreditOperationType,
			amount:        zeroAmount,
			wantErrs:      []error{domain.ErrTransactionInvalidTransactionAmount},
		},
		{
			name:          "negative amount returns sign error",
			accountID:     1,
			operationType: validCreditOperationType,
			amount:        negativeAmount,
			wantErrs:      []error{domain.ErrTransactionNegativeAmount},
		},
		{
			name:          "multiple validation errors are returned",
			accountID:     0,
			operationType: domain.OperationType{}, //nolint: exhaustruct
			amount:        zeroAmount,
			wantErrs: []error{
				domain.ErrTransactionInvalidAccountID,
				domain.ErrTransactionInvalidTransactionAmount,
				domain.ErrTransactionInvalidOperationTypeID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			before := time.Now()
			got, gotErrs := domain.NewTransaction(tt.accountID, tt.operationType, tt.amount)
			after := time.Now()

			if len(tt.wantErrs) > 0 {
				require.Nil(t, got)
				require.Len(t, gotErrs, len(tt.wantErrs))
				for i, wantErr := range tt.wantErrs {
					assert.ErrorIs(t, gotErrs[i], wantErr)
				}
				return
			}

			require.Nil(t, gotErrs)
			require.NotNil(t, got)

			assert.Equal(t, int64(0), got.ID)
			assert.Equal(t, tt.accountID, got.AccountID)
			assert.Equal(t, tt.operationType, got.OperationType)
			assert.Equal(t, domain.BRL, got.Currency)
			assert.WithinRange(t, got.EventDate, before, after)

			expectedAmount := tt.amount
			if tt.operationType.IsDebit() {
				expectedAmount = tt.amount.Neg()
			}
			assert.Equal(t, expectedAmount.Int64(), got.Amount.Int64())

			assert.Equal(t, got.Amount.IsNegative(), got.IsDebit())
			assert.Equal(t, got.Amount.IsPositive(), got.IsCredit())
		})
	}
}

func TestTransaction_IsDebit(t *testing.T) {
	tests := []struct {
		name   string
		amount string
		want   bool
	}{
		{name: "negative amount is debit", amount: "-10.00", want: true},
		{name: "positive amount is not debit", amount: "10.00", want: false},
		{name: "zero amount is not debit", amount: "0", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			amount, err := domain.NewBRLMoneyFromString(tt.amount)
			require.NoError(t, err)

			tr := domain.Transaction{Amount: amount} //nolint: exhaustruct
			got := tr.IsDebit()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTransaction_IsCredit(t *testing.T) {
	tests := []struct {
		name   string
		amount string
		want   bool
	}{
		{name: "positive amount is credit", amount: "10.00", want: true},
		{name: "negative amount is not credit", amount: "-10.00", want: false},
		{name: "zero amount is not credit", amount: "0", want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			amount, err := domain.NewBRLMoneyFromString(tt.amount)
			require.NoError(t, err)

			tr := domain.Transaction{Amount: amount} //nolint: exhaustruct
			got := tr.IsCredit()
			assert.Equal(t, tt.want, got)
		})
	}
}
