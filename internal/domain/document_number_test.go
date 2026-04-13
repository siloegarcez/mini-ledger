package domain_test

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDocumentNumber(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    string
		wantErr error
	}{
		{
			name:    "valid document number",
			value:   "abc123",
			want:    "abc123",
			wantErr: nil,
		},
		{
			name:    "lowercases input",
			value:   "ABC123XYZ",
			want:    "abc123xyz",
			wantErr: nil,
		},
		{
			name:    "removes all whitespace",
			value:   " A B\tC\n1 2 3 ",
			want:    "abc123",
			wantErr: nil,
		},
		{
			name:    "max length is accepted",
			value:   "ABCDEFGHIJKLMNO",
			want:    "abcdefghijklmno",
			wantErr: nil,
		},
		{
			name:    "greater than max length returns invalid length",
			value:   "abcdefghijklmnop",
			want:    "",
			wantErr: domain.ErrDocumentNumberInvalidLen,
		},
		{
			name:    "empty string returns empty error",
			value:   "",
			want:    "",
			wantErr: domain.ErrDocumentNumberEmpty,
		},
		{
			name:    "only whitespace returns empty error",
			value:   " \t\n\r ",
			want:    "",
			wantErr: domain.ErrDocumentNumberEmpty,
		},
		{
			name:    "non breaking spaces are removed",
			value:   "AB\u00A0123",
			want:    "ab123",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := domain.NewDocumentNumber(tt.value)

			if tt.wantErr != nil {
				require.Error(t, gotErr)
				require.ErrorIs(t, gotErr, tt.wantErr)
				return
			}

			require.NoError(t, gotErr)
			assert.Equal(t, tt.want, got.String())
		})
	}
}
