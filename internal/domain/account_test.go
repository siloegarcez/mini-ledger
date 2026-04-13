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
		wantErr        error
	}{
		{
			name:           "empty document number returns validation error",
			documentNumber: domain.DocumentNumber{},
			wantErr:        domain.ErrAccountDocumentNumberEmpty,
		},
		{
			name:           "valid document number creates account",
			documentNumber: validDocumentNumber,
			wantErr:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := domain.NewAccount(tt.documentNumber)

			if tt.wantErr != nil {
				require.ErrorIs(t, gotErr, tt.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, gotErr)

			require.NotNil(t, got)

			assert.Equal(t, int64(0), got.ID)

			assert.Equal(t, tt.documentNumber.String(), got.DocumentNumber.String())

			assert.True(t, got.CreatedAt.IsZero())

			assert.True(t, got.UpdatedAt.IsZero())

			assert.True(t, got.CreatedAt.Equal(got.UpdatedAt))
		})
	}
}
