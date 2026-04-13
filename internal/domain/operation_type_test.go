package domain_test

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOperationType(t *testing.T) {
	tests := []struct {
		name            string
		operationTypeID int64
		description     string
		signMultiplier  int16
		want            *domain.OperationType
		wantErrs        []error
	}{
		{
			name:            "valid credit operation type",
			operationTypeID: 1,
			description:     "PAYMENT",
			signMultiplier:  domain.CreditSignMultiplier,
			want: &domain.OperationType{
				OperationTypeID: 1,
				Description:     "PAYMENT",
				SignMultiplier:  domain.CreditSignMultiplier,
			},
			wantErrs: nil,
		},
		{
			name:            "valid debit operation type",
			operationTypeID: 2,
			description:     "INSTALLMENT PURCHASE",
			signMultiplier:  domain.DebitSignMultiplier,
			want: &domain.OperationType{
				OperationTypeID: 2,
				Description:     "INSTALLMENT PURCHASE",
				SignMultiplier:  domain.DebitSignMultiplier,
			},
			wantErrs: nil,
		},
		{
			name:            "invalid operation type id",
			operationTypeID: 0,
			description:     "PAYMENT",
			signMultiplier:  domain.CreditSignMultiplier,
			want:            nil,
			wantErrs:        []error{domain.ErrOperationTypeInvalidOperationTypeID},
		},
		{
			name:            "invalid sign multiplier",
			operationTypeID: 1,
			description:     "PURCHASE",
			signMultiplier:  0,
			want:            nil,
			wantErrs:        []error{domain.ErrOperationTypeInvalidSignMultiplier},
		},
		{
			name:            "empty description",
			operationTypeID: 1,
			description:     "",
			signMultiplier:  domain.CreditSignMultiplier,
			want:            nil,
			wantErrs:        []error{domain.ErrOperationTypeEmptyDescription},
		},
		{
			name:            "blank description",
			operationTypeID: 1,
			description:     "   ",
			signMultiplier:  domain.CreditSignMultiplier,
			want:            nil,
			wantErrs:        []error{domain.ErrOperationTypeEmptyDescription},
		},
		{
			name:            "description too long",
			operationTypeID: 1,
			description:     "1234567890123456789012345678901",
			signMultiplier:  domain.CreditSignMultiplier,
			want:            nil,
			wantErrs:        []error{domain.ErrOperationTypeInvalidDescriptionLen},
		},
		{
			name:            "multiple validation errors",
			operationTypeID: -10,
			description:     "",
			signMultiplier:  99,
			want:            nil,
			wantErrs: []error{
				domain.ErrOperationTypeInvalidOperationTypeID,
				domain.ErrOperationTypeInvalidSignMultiplier,
				domain.ErrOperationTypeEmptyDescription,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErrs := domain.NewOperationType(
				tt.operationTypeID,
				tt.description,
				tt.signMultiplier,
			)

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
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOperationType_IsDebit(t *testing.T) {
	tests := []struct {
		name           string
		signMultiplier int16
		want           bool
	}{
		{
			name:           "debit multiplier returns true",
			signMultiplier: domain.DebitSignMultiplier,
			want:           true,
		},
		{
			name:           "credit multiplier returns false",
			signMultiplier: domain.CreditSignMultiplier,
			want:           false,
		},
		{
			name:           "invalid multiplier returns false",
			signMultiplier: 0,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ot := domain.OperationType{SignMultiplier: tt.signMultiplier} //nolint: exhaustruct
			got := ot.IsDebit()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOperationType_IsCredit(t *testing.T) {
	tests := []struct {
		name           string
		signMultiplier int16
		want           bool
	}{
		{
			name:           "credit multiplier returns true",
			signMultiplier: domain.CreditSignMultiplier,
			want:           true,
		},
		{
			name:           "debit multiplier returns false",
			signMultiplier: domain.DebitSignMultiplier,
			want:           false,
		},
		{
			name:           "invalid multiplier returns false",
			signMultiplier: 0,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ot := domain.OperationType{SignMultiplier: tt.signMultiplier} //nolint: exhaustruct
			got := ot.IsCredit()

			assert.Equal(t, tt.want, got)
		})
	}
}
