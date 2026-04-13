package postgres

import (
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mapTransactionModelToDomain(t *testing.T) {
	tests := []struct {
		name        string
		transaction *model.Transactions
		want        *domain.Transaction
		wantErr     bool
	}{
		{
			name: "maps model to domain successfully",
			transaction: &model.Transactions{
				ID:              10,
				AccountID:       20,
				OperationTypeID: 4,
				Amount:          12345,
				Currency:        "BRL",
				Scale:           domain.BRLScale,
				EventDate:       time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
			},
			want: &domain.Transaction{
				ID:              10,
				AccountID:       20,
				OperationTypeID: 4,
				Currency:        "BRL",
				Amount:          mustMoneyFromInt64(t, 12345, domain.BRLScale),
				EventDate:       time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "invalid scale returns error",
			transaction: &model.Transactions{ // nolint: exhaustruct
				Amount: 100,
				Scale:  1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := mapTransactionModelToDomain(tt.transaction)
			if gotErr != nil {
				assert.True(t, tt.wantErr)
				return
			}
			require.False(t, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mapTransactionDomainToModel(t *testing.T) {
	tests := []struct {
		name        string
		transaction *domain.Transaction
		want        *model.Transactions
	}{
		{
			name: "maps domain to model",
			transaction: &domain.Transaction{
				ID:              11,
				AccountID:       22,
				OperationTypeID: 3,
				Amount:          mustMoneyFromInt64(t, -5010, domain.BRLScale),
				Currency:        "BRL",
				EventDate:       time.Date(2026, time.April, 11, 8, 30, 0, 0, time.UTC),
			},
			want: &model.Transactions{
				ID:              11,
				AccountID:       22,
				OperationTypeID: 3,
				Amount:          -5010,
				Scale:           domain.BRLScale,
				Currency:        "BRL",
				EventDate:       time.Date(2026, time.April, 11, 8, 30, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapTransactionDomainToModel(tt.transaction)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustMoneyFromInt64(t *testing.T, amount int64, scale int16) domain.BRLMoney {
	t.Helper()
	m, err := domain.NewMoneyFromInt64(amount, scale)
	require.NoError(t, err)
	return m
}
