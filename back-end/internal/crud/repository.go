package crud

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"
)

type Repository[T any] struct {
	res *Resource[T]
}

func NewRepository[T any](res *Resource[T]) *Repository[T] {
	return &Repository[T]{res: res}
}

func (r *Repository[T]) Search(ctx context.Context, sc queries.SearchCriteria) ([]T, error) {
	rows, err := database.PerformSearch(ctx, sc, r.res.Columns, r.res.Table, r.res.Search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []T{}

	for rows.Next() {
		x, err := r.res.Scan(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository[T]) Get(ctx context.Context, id any) (*T, error) {
	row, err := database.PerformGet(ctx, id, r.res.Columns, r.res.Table, r.res.Search)
	if err != nil {
		return nil, err
	}

	x, err := r.res.Scan(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("%s '%s' not found", r.res.Name, id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
