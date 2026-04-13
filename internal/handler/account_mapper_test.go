package handler

import (
	"mini-ledger/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_mapDomainAccountToAccountCreateResponse(t *testing.T) {
	tests := []struct {
		name string
		acc  *domain.Account
		want AccountCreateResponse
	}{
		{
			name: "maps account to create response",
			acc: &domain.Account{ // nolint: exhaustruct
				ID:             10,
				DocumentNumber: mustDocNumber(t, "12345678901"),
			},
			want: AccountCreateResponse{
				AccountID:      10,
				DocumentNumber: "12345678901",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapDomainAccountToAccountCreateResponse(tt.acc)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_mapDomainAccountToAccountGetByIDResponse(t *testing.T) {
	tests := []struct {
		name string
		acc  *domain.Account
		want AccountGetByIDResponse
	}{
		{
			name: "maps account to get by id response",
			acc: &domain.Account{ // nolint: exhaustruct
				ID:             11,
				DocumentNumber: mustDocNumber(t, "abc123"),
			},
			want: AccountGetByIDResponse{
				AccountID:      11,
				DocumentNumber: "abc123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapDomainAccountToAccountGetByIDResponse(tt.acc)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustDocNumber(t *testing.T, value string) domain.DocumentNumber {
	t.Helper()
	d, err := domain.NewDocumentNumber(value)
	require.NoError(t, err)
	return d
}
