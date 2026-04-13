package postgres

import (
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mapAccountsModelToDomain(t *testing.T) {
	tests := []struct {
		name    string
		account *model.Accounts
		want    *domain.Account
		wantErr bool
	}{
		{
			name: "maps account model to domain successfully",
			account: &model.Accounts{
				ID:             1,
				DocumentNumber: "abc123",
				CreatedAt:      time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2026, time.April, 10, 13, 0, 0, 0, time.UTC),
			},
			want: &domain.Account{
				ID:             1,
				DocumentNumber: mustDocumentNumber(t, "abc123"),
				CreatedAt:      time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2026, time.April, 10, 13, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "invalid document number returns error",
			account: &model.Accounts{ // nolint: exhaustruct
				ID:             1,
				DocumentNumber: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := mapAccountsModelToDomain(tt.account)
			if gotErr != nil {
				assert.True(t, tt.wantErr)
				return
			}
			require.False(t, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mapDomainAccountToModel(t *testing.T) {
	tests := []struct {
		name    string
		account *domain.Account
		want    *model.Accounts
	}{
		{
			name: "maps domain account to model",
			account: &domain.Account{
				ID:             12,
				DocumentNumber: mustDocumentNumber(t, "xyz987"),
				CreatedAt:      time.Date(2026, time.April, 12, 8, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2026, time.April, 12, 9, 0, 0, 0, time.UTC),
			},
			want: &model.Accounts{
				ID:             12,
				DocumentNumber: "xyz987",
				CreatedAt:      time.Date(2026, time.April, 12, 8, 0, 0, 0, time.UTC),
				UpdatedAt:      time.Date(2026, time.April, 12, 9, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapDomainAccountToModel(tt.account)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustDocumentNumber(t *testing.T, value string) domain.DocumentNumber {
	t.Helper()
	d, err := domain.NewDocumentNumber(value)
	require.NoError(t, err)
	return d
}
