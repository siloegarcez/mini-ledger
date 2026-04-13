package postgres

import (
	"mini-ledger/gen/dev/public/model"
	"mini-ledger/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_mapOperationTypeModelToDomain(t *testing.T) {
	tests := []struct {
		name          string
		operationType *model.OperationsTypes
		want          *domain.OperationType
	}{
		{
			name: "maps operation type model to domain",
			operationType: &model.OperationsTypes{
				OperationTypeID: 4,
				Description:     "PAYMENT",
				SignMultiplier:  domain.CreditSignMultiplier,
				CreatedAt:       time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
			},
			want: &domain.OperationType{
				OperationTypeID: 4,
				Description:     "PAYMENT",
				SignMultiplier:  domain.CreditSignMultiplier,
				CreatedAt:       time.Date(2026, time.April, 10, 12, 0, 0, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOperationTypeModelToDomain(tt.operationType)
			assert.Equal(t, tt.want, got)
		})
	}
}
