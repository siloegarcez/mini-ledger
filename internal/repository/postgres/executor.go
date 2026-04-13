package postgres

import (
	"context"
	"database/sql"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

func NewExecutor(db *sql.DB) Executor {
	return &executor{
		db: db,
	}
}

type executor struct {
	db *sql.DB
}

// ExecContext implements [Executor].
func (e *executor) ExecContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (sql.Result, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.ExecContext(ctx, query, args...)
	}
	return e.db.ExecContext(ctx, query, args...)
}

// QueryContext implements [Executor].
func (e *executor) QueryContext(
	ctx context.Context,
	query string,
	args ...interface{},
) (*sql.Rows, error) {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}
	return e.db.QueryContext(ctx, query, args...)
}
