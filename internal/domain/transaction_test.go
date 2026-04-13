package domain_test

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTransaction(t *testing.T) {
	tests := []struct {
		name            string
		accountID       int64
		operationTypeID int64
		amountCents     int64
		wantErr         error
	}{
		{
			name:            "invalid account id returns error",
			accountID:       0,
			operationTypeID: 1,
			amountCents:     1000,
			wantErr:         domain.ErrTransactionInvalidAccountID,
		},
		{
			name:            "zero amount returns error",
			accountID:       1,
			operationTypeID: 1,
			amountCents:     0,
			wantErr:         domain.ErrTransactionInvalidTransactionAmount,
		},
		{
			name:            "invalid operation type id returns error",
			accountID:       1,
			operationTypeID: 0,
			amountCents:     1000,
			wantErr:         domain.ErrTransactionInvalidOperationTypeID,
		},
		{
			name:            "negative amount returns error",
			accountID:       1,
			operationTypeID: 1,
			amountCents:     -1000,
			wantErr:         domain.ErrTransactionNegativeAmount,
		},
		{
			name:            "returns first validation error",
			accountID:       -1,
			operationTypeID: 0,
			amountCents:     0,
			wantErr:         domain.ErrTransactionInvalidAccountID,
		},
		{
			name:            "valid transaction",
			accountID:       1,
			operationTypeID: 4,
			amountCents:     1050,
			wantErr:         nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			amount, err := domain.NewMoneyFromInt64(tt.amountCents, domain.BRLScale)
			require.NoError(t, err)

			got, gotErr := domain.NewTransaction(tt.accountID, tt.operationTypeID, amount)

			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, gotErr)
			require.NotNil(t, got)
			assert.Equal(t, tt.accountID, got.AccountID)
			assert.Equal(t, tt.operationTypeID, got.OperationTypeID)
			assert.Equal(t, amount.Int64(), got.Amount.Int64())
			assert.Equal(t, int16(domain.BRLScale), got.Amount.Scale())
			assert.Equal(t, domain.BRL, got.Currency)
			assert.Zero(t, got.EventDate)
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
