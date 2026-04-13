package handler

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mapDomainTransactionToTransactionCreateResponse(t *testing.T) {
	tests := []struct {
		name        string
		transaction *domain.Transaction
		want        TransactionCreateResponse
	}{
		{
			name: "maps transaction to create response",
			transaction: &domain.Transaction{ // nolint: exhaustruct
				ID:              10,
				AccountID:       20,
				OperationTypeID: 4,
				Amount:          mustMoneyFromString(t, "123.45"),
			},
			want: TransactionCreateResponse{
				TransactionID:   10,
				AccountID:       20,
				OperationTypeID: 4,
				Amount:          Number("123.45"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapDomainTransactionToTransactionCreateResponse(tt.transaction)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustMoneyFromString(t *testing.T, value string) domain.BRLMoney {
	t.Helper()
	m, err := domain.NewBRLMoneyFromString(value)
	require.NoError(t, err)
	return m
}
