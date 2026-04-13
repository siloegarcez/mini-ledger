package dberrs

import (
	"errors"
	"mini-ledger/internal/database"

	"github.com/jackc/pgx/v5/pgconn"
)

type postgresError struct {
	msg string
}

func (e postgresError) Error() string {
	return e.msg
}

func newErr(msg string) postgresError {
	return postgresError{msg: msg}
}

var (
	ErrForeignKeyViolation   = newErr("foreign key violation")
	ErrDuplicateKeyViolation = newErr("duplicate key violation")
	ErrNotImplemented        = newErr("not implemented")
	ErrNotPgError            = newErr("not a pgconn.PgError")
)

func getPostgresSentinelErr(err error) (*pgconn.PgError, error) {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return nil, ErrNotPgError
	}

	var sentinelErr error
	sentinelErr = ErrNotImplemented

	switch pgErr.SQLState() {
	case "23503":
		sentinelErr = ErrForeignKeyViolation
	case "23505":
		sentinelErr = ErrDuplicateKeyViolation
	}

	return pgErr, sentinelErr
}

func Is(err error, target postgresError) bool {
	if err == nil {
		return false
	}
	_, sentinelErr := getPostgresSentinelErr(err)
	return errors.Is(sentinelErr, target)
}

func IsWithConstraint(err error, target postgresError, constraint database.Constraint) bool {
	if err == nil {
		return false
	}
	pgErr, sentinelErr := getPostgresSentinelErr(err)
	return errors.Is(sentinelErr, target) && pgErr.ConstraintName == string(constraint)
}

func Unwrap(err error) *pgconn.PgError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr
	}
	return nil
}
