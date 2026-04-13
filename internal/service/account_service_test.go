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

type stubAccountRepository struct {
	createFn  func(ctx context.Context, account *domain.Account) (*domain.Account, error)
	getByIDFn func(ctx context.Context, id int64) (*domain.Account, error)
}

func (s stubAccountRepository) Create(
	ctx context.Context,
	account *domain.Account,
) (*domain.Account, error) {
	return s.createFn(ctx, account)
}

func (s stubAccountRepository) GetByID(ctx context.Context, id int64) (*domain.Account, error) {
	return s.getByIDFn(ctx, id)
}

type noopTransactor struct{}

func (n noopTransactor) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	return fn(ctx)
}

func TestAccountService_Create(t *testing.T) {
	tests := []struct {
		name       string
		doc        string
		repoCreate func(ctx context.Context, account *domain.Account) (*domain.Account, error)
		assertErr  func(t *testing.T, err error)
		assertAcc  func(t *testing.T, acc *domain.Account)
	}{
		{
			name: "invalid document number returns unprocessable entity",
			doc:  "",
			repoCreate: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
				t.Fatalf("repository must not be called for invalid input")
				return nil, nil //nolint:nilnil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assertStatusError(t, err, http.StatusUnprocessableEntity)
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				assert.Nil(t, acc)
			},
		},
		{
			name: "duplicate key violation maps to conflict",
			doc:  "12345678900",
			repoCreate: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
				return nil, &pgconn.PgError{Code: "23505"} //nolint:exhaustruct
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assertStatusError(t, err, http.StatusConflict)
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				assert.Nil(t, acc)
			},
		},
		{
			name: "repository error is returned as is",
			doc:  "12345678900",
			repoCreate: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
				return nil, errors.New("db offline") //nolint:err113
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assert.EqualError(t, err, "db offline")
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				assert.Nil(t, acc)
			},
		},
		{
			name: "successful create returns created account",
			doc:  "ABC123",
			repoCreate: func(ctx context.Context, account *domain.Account) (*domain.Account, error) {
				assert.Equal(t, "abc123", account.DocumentNumber.String())
				return &domain.Account{ //nolint:exhaustruct
					ID:             10,
					DocumentNumber: account.DocumentNumber,
				}, nil
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				require.NotNil(t, acc)
				assert.Equal(t, int64(10), acc.ID)
				assert.Equal(t, "abc123", acc.DocumentNumber.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewAccountService(
				stubAccountRepository{createFn: tt.repoCreate}, //nolint:exhaustruct
				noopTransactor{},
			)
			acc, err := svc.Create(t.Context(), tt.doc)

			tt.assertErr(t, err)
			tt.assertAcc(t, acc)
		})
	}
}

func TestAccountService_GetByID(t *testing.T) {
	tests := []struct {
		name      string
		repoGet   func(ctx context.Context, id int64) (*domain.Account, error)
		assertErr func(t *testing.T, err error)
		assertAcc func(t *testing.T, acc *domain.Account)
	}{
		{
			name: "not found maps to 404",
			repoGet: func(ctx context.Context, id int64) (*domain.Account, error) {
				return nil, qrm.ErrNoRows
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assertStatusError(t, err, http.StatusNotFound)
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				assert.Nil(t, acc)
			},
		},
		{
			name: "repository error is returned as is",
			repoGet: func(ctx context.Context, id int64) (*domain.Account, error) {
				return nil, errors.New("storage unavailable") //nolint:err113
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				require.Error(t, err)
				assert.EqualError(t, err, "storage unavailable")
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				assert.Nil(t, acc)
			},
		},
		{
			name: "success returns account",
			repoGet: func(ctx context.Context, id int64) (*domain.Account, error) {
				doc, err := domain.NewDocumentNumber("123")
				require.NoError(t, err)
				return &domain.Account{ID: id, DocumentNumber: doc}, nil //nolint:exhaustruct
			},
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				assert.NoError(t, err)
			},
			assertAcc: func(t *testing.T, acc *domain.Account) {
				t.Helper()
				require.NotNil(t, acc)
				assert.Equal(t, int64(33), acc.ID)
				assert.Equal(t, "123", acc.DocumentNumber.String())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewAccountService(
				stubAccountRepository{getByIDFn: tt.repoGet}, //nolint:exhaustruct
				noopTransactor{},
			)
			acc, err := svc.GetByID(t.Context(), 33)

			tt.assertErr(t, err)
			tt.assertAcc(t, acc)
		})
	}
}

func assertStatusError(t *testing.T, err error, status int) {
	t.Helper()

	var statusErr huma.StatusError
	require.ErrorAs(t, err, &statusErr)
	assert.Equal(t, status, statusErr.GetStatus())
}
