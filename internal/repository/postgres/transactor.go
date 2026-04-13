package postgres

import (
	"context"
	"database/sql"
	"errors"
	"mini-ledger/internal/domain"

	"fmt"

	"github.com/rs/zerolog/log"
)

type transactor struct {
	db *sql.DB
}

var (
	ErrTransactionFailed = errors.New("transaction failed")
	ErrRollbackFailed    = errors.New("rollback transaction failed")
	ErrCommitFailed      = errors.New("commit transaction failed")
)

func NewTransactor(db *sql.DB) domain.Transactor {
	return &transactor{db: db}
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	if tx == nil {
		return ctx
	}
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (t *transactor) WithinTransaction(
	ctx context.Context,
	tFunc func(ctx context.Context) error,
) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Err(err).Msg("Failed to begin transaction")
		return fmt.Errorf("%w: %w", ErrTransactionFailed, err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Ctx(ctx).
					Error().
					Err(rollbackErr).
					Msg("Failed to rollback transaction after panic")
			}
			log.Ctx(ctx).Error().Interface("panic", p).Msg("Panic during transaction")
			panic(p) // Re-throw panic
		}
	}()

	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		// Rollback on error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Ctx(ctx).Error().Err(rollbackErr).Msg("Failed to rollback transaction")
			return fmt.Errorf("%w: %w (original error: %w)", ErrRollbackFailed, rollbackErr, err)
		}
		log.Ctx(ctx).Debug().Err(err).Msg("Transaction rolled back")
		return err
	}

	// Commit transaction
	if commitErr := tx.Commit(); commitErr != nil {
		log.Ctx(ctx).Error().Err(commitErr).Msg("Failed to commit transaction")
		return fmt.Errorf("%w: %w", ErrCommitFailed, commitErr)
	}

	log.Debug().Msg("Transaction committed successfully")
	return nil
}
