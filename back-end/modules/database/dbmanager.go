package database

import (
	"context"
	"database/sql"
)

type contextKey string

const txKey contextKey = "transaction_key"

func InjectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func ExtractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}
