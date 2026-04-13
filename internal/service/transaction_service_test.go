package service

import (
	"context"
	"errors"
	"mini-ledger/internal/domain"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubTransactionRepository struct {
	createFn func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
}

func (s stubTransactionRepository) Create(
	ctx context.Context,
	transaction *domain.Transaction,
) (*domain.Transaction, error) {
	return s.createFn(ctx, transaction)
}

type stubOperationTypeRepository struct {
	getByIDFn func(ctx context.Context, id int64) (*domain.OperationType, error)
}

func (s stubOperationTypeRepository) GetByID(
	ctx context.Context,
	id int64,
) (*domain.OperationType, error) {
	return s.getByIDFn(ctx, id)
}

func TestTransactionService_Create(t *testing.T) {
	tests := []struct {
		name        string
		cmd         *CreateTransactionCommand
		opTypeGet   func(ctx context.Context, id int64) (*domain.OperationType, error)
		repoCreate  func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error)
		assertErr   func(t *testing.T, err error)
		assertTrans func(t *testing.T, tr *domain.Transaction)
	}{
		{
			name: "operation type not found maps to bad request",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       1,
				OperationTypeID: 999,
				Amount:          "10.00",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return nil, qrm.ErrNoRows
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				t.Fatalf("repository must not be called when operation type is invalid")
				return nil, nil //nolint:nilnil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assertStatusErrorTx(t, err, http.StatusBadRequest)
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				assert.Nil(t, tr)
			},
		},
		{
			name: "invalid amount maps to unprocessable entity",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          "abc",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return &domain.OperationType{ //nolint:exhaustruct
					OperationTypeID: id,
					Description:     "PAYMENT",
					SignMultiplier:  domain.CreditSignMultiplier,
				}, nil
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				t.Fatalf("repository must not be called when amount is invalid")
				return nil, nil //nolint: nilnil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assertStatusErrorTx(t, err, http.StatusUnprocessableEntity)
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				assert.Nil(t, tr)
			},
		},
		{
			name: "foreign key violation maps to invalid account id",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       7,
				OperationTypeID: 4,
				Amount:          "12.34",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return &domain.OperationType{ //nolint:exhaustruct
					OperationTypeID: id,
					Description:     "PAYMENT",
					SignMultiplier:  domain.CreditSignMultiplier,
				}, nil
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				return nil, &pgconn.PgError{Code: "23503"} //nolint:exhaustruct
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assertStatusErrorTx(t, err, http.StatusBadRequest)
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				assert.Nil(t, tr)
			},
		},
		{
			name: "debit operation stores negative amount",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       9,
				OperationTypeID: 1,
				Amount:          "10.00",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return &domain.OperationType{ //nolint:exhaustruct
					OperationTypeID: id,
					Description:     "PURCHASE",
					SignMultiplier:  domain.DebitSignMultiplier,
				}, nil
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				assert.Equal(t, int64(-1000), transaction.Amount.Int64())
				return &domain.Transaction{ //nolint:exhaustruct
					ID:              55,
					AccountID:       transaction.AccountID,
					OperationTypeID: transaction.OperationTypeID,
					Amount:          transaction.Amount,
					Currency:        transaction.Currency,
				}, nil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				require.NotNil(t, tr)
				assert.Equal(t, int64(55), tr.ID)
				assert.Equal(t, int64(-1000), tr.Amount.Int64())
			},
		},
		{
			name: "credit operation keeps positive amount",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       11,
				OperationTypeID: 4,
				Amount:          "10.00",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return &domain.OperationType{ //nolint:exhaustruct
					OperationTypeID: id,
					Description:     "PAYMENT",
					SignMultiplier:  domain.CreditSignMultiplier,
				}, nil
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				assert.Equal(t, int64(1000), transaction.Amount.Int64())
				return &domain.Transaction{ //nolint:exhaustruct
					ID:              77,
					AccountID:       transaction.AccountID,
					OperationTypeID: transaction.OperationTypeID,
					Amount:          transaction.Amount,
					Currency:        transaction.Currency,
				}, nil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				require.NotNil(t, tr)
				assert.Equal(t, int64(77), tr.ID)
				assert.Equal(t, int64(1000), tr.Amount.Int64())
			},
		},
		{
			name: "operation type repository error is propagated",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          "10.00",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return nil, errors.New("operation type db error") //nolint:err113
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				t.Fatalf("repository must not be called when operation type query fails")
				return nil, nil //nolint:nilnil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assert.EqualError(t, err, "operation type db error")
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				assert.Nil(t, tr)
			},
		},
		{
			name: "transaction repository generic error is propagated",
			cmd: &CreateTransactionCommand{ //nolint:exhaustruct
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          "10.00",
			},
			opTypeGet: func(ctx context.Context, id int64) (*domain.OperationType, error) {
				return &domain.OperationType{ //nolint:exhaustruct
					OperationTypeID: id,
					Description:     "PAYMENT",
					SignMultiplier:  domain.CreditSignMultiplier,
				}, nil
			},
			repoCreate: func(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
				return nil, errors.New("insert failed") //nolint:err113
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assert.EqualError(t, err, "insert failed")
			},
			assertTrans: func(t *testing.T, tr *domain.Transaction) {
				t.Helper()
				assert.Nil(t, tr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewTransactionService(
				stubTransactionRepository{createFn: tt.repoCreate},
				stubOperationTypeRepository{getByIDFn: tt.opTypeGet},
				noopTransactor{},
			)

			tr, err := svc.Create(t.Context(), tt.cmd)

			tt.assertErr(t, err)
			tt.assertTrans(t, tr)
		})
	}
}

func assertStatusErrorTx(t *testing.T, err error, status int) {
	t.Helper()

	var statusErr huma.StatusError
	require.ErrorAs(t, err, &statusErr)
	assert.Equal(t, status, statusErr.GetStatus())
}
