package domain_test

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	validDocumentNumber, err := domain.NewDocumentNumber("12345678901")
	require.NoError(t, err)

	tests := []struct {
		name           string
		documentNumber domain.DocumentNumber
		wantErrs       []error
	}{
		{
			name:           "empty document number returns validation error",
			documentNumber: domain.DocumentNumber{},
			wantErrs:       []error{domain.ErrAccountDocumentNumberEmpty},
		},
		{
			name:           "valid document number creates account",
			documentNumber: validDocumentNumber,
			wantErrs:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErrs := domain.NewAccount(tt.documentNumber)

			require.Len(t, gotErrs, len(tt.wantErrs))

			for i := range tt.wantErrs {
				require.ErrorIs(t, gotErrs[i], tt.wantErrs[i])
			}

			if len(tt.wantErrs) > 0 {
				assert.Nil(t, got)
				return
			}

			require.NotNil(t, got)

			assert.Equal(t, int64(0), got.ID)

			assert.Equal(t, tt.documentNumber.String(), got.DocumentNumber.String())

			assert.False(t, got.CreatedAt.IsZero())

			assert.False(t, got.UpdatedAt.IsZero())

			assert.True(t, got.CreatedAt.Equal(got.UpdatedAt))
		})
	}
}
