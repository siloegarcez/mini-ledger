package dberrs

import (
	"errors"
	"mini-ledger/internal/database"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIs(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target postgresError
		want   bool
	}{
		{
			name:   "nil error returns false",
			err:    nil,
			target: ErrForeignKeyViolation,
			want:   false,
		},
		{
			name:   "non postgres error returns false",
			err:    errors.New("boom"), //nolint
			target: ErrForeignKeyViolation,
			want:   false,
		},
		{
			name: "postgres foreign key violation returns true",
			err: &pgconn.PgError{ //nolint:exhaustruct
				Code: "23503",
			},
			target: ErrForeignKeyViolation,
			want:   true,
		},
		{
			name: "postgres duplicate key does not match foreign key target",
			err: &pgconn.PgError{ //nolint:exhaustruct
				Code: "23505",
			},
			target: ErrForeignKeyViolation,
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Is(tt.err, tt.target)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsWithConstraint(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		target     postgresError
		constraint database.Constraint
		want       bool
	}{
		{
			name:       "nil error returns false",
			err:        nil,
			target:     ErrForeignKeyViolation,
			constraint: database.Constraint("fk_transactions_account_id"),
			want:       false,
		},
		{
			name:       "non postgres error returns false",
			err:        errors.New("boom"), //nolint
			target:     ErrForeignKeyViolation,
			constraint: database.Constraint("fk_transactions_account_id"),
			want:       false,
		},
		{
			name: "matching postgres code and constraint returns true",
			err: &pgconn.PgError{ //nolint:exhaustruct
				Code:           "23503",
				ConstraintName: "fk_transactions_account_id",
			},
			target:     ErrForeignKeyViolation,
			constraint: database.Constraint("fk_transactions_account_id"),
			want:       true,
		},
		{
			name: "matching code but different constraint returns false",
			err: &pgconn.PgError{ //nolint:exhaustruct
				Code:           "23503",
				ConstraintName: "fk_transactions_operation_type_id",
			},
			target:     ErrForeignKeyViolation,
			constraint: database.Constraint("fk_transactions_account_id"),
			want:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsWithConstraint(tt.err, tt.target, tt.constraint)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUnwrap(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *pgconn.PgError
	}{
		{
			name: "postgres error unwraps successfully",
			err: &pgconn.PgError{ //nolint:exhaustruct
				Code:           "23505",
				ConstraintName: "accounts_document_number_key",
			},
			want: &pgconn.PgError{ //nolint:exhaustruct
				Code:           "23505",
				ConstraintName: "accounts_document_number_key",
			},
		},
		{
			name: "non postgres error returns nil",
			err:  errors.New("boom"), // nolint
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unwrap(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}
