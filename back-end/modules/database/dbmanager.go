package database

import (
	"backend/utils/queries"
	"context"
	"database/sql"
	"fmt"
)

type contextKey string

const txKey contextKey = "transaction_key"

// Scanner is satisfied by both *sql.Rows and *sql.Row, so one scan function serves search and get by id.
type Scanner interface {
	Scan(dest ...any) error
}

func InjectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func ExtractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// WithReadTx runs fn inside a read-only transaction injected into its context, committing on success.
func WithReadTx[T any](ctx context.Context, db *sql.DB, fn func(ctx context.Context) (T, error)) (T, error) {
	var zero T

	tx, err := db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return zero, err
	}
	defer tx.Rollback()

	result, err := fn(InjectTx(ctx, tx))
	if err != nil {
		return zero, err
	}

	if err := tx.Commit(); err != nil {
		return zero, err
	}
	return result, nil
}

func PerformSearch(ctx context.Context, sc queries.SearchCriteria, cols []string, table string, cfg queries.EntitySearchConfig) (*sql.Rows, error) {
	exec := ExtractTx(ctx)
	if exec == nil {
		return nil, fmt.Errorf("No transaction provided")
	}

	query, args, err := queries.BuildSearchQuery(cols, table, cfg, sc)

	if err != nil {
		return nil, err
	}

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// PerformGet selects one row by id. Scanning it yields sql.ErrNoRows when it does not exist.
func PerformGet(ctx context.Context, id any, cols []string, table string, cfg queries.EntitySearchConfig) (*sql.Row, error) {
	exec := ExtractTx(ctx)
	if exec == nil {
		return nil, fmt.Errorf("No transaction provided")
	}

	return exec.QueryRowContext(ctx, queries.BuildGetQuery(cols, table, cfg), id), nil
}
