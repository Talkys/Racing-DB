package database

import (
	"backend/utils/queries"
	"context"
	"database/sql"
	"fmt"
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

func PerformSearch(ctx context.Context, sc queries.SearchCriteria, cols []string, table string, cfg queries.EntitySearchConfig) (*sql.Rows, error) {
	exec := ExtractTx(ctx)
	if exec == nil {
		return nil, fmt.Errorf("No transaction provided")
	}

	query, args, err := queries.BuildSearchQuery(cols, table, cfg, sc)

	if err != nil {
		return nil, err
	}

	rows, err := exec.QueryContext(ctx, query, args[0])
	if err != nil {
		return nil, err
	}

	return rows, nil
}
