package crud

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type Service[T any] struct {
	db   *sql.DB
	core *Core[T]
}

func NewService[T any](db *sql.DB, res *Resource[T]) *Service[T] {
	return &Service[T]{
		db:   db,
		core: NewCore(res),
	}
}

func (s *Service[T]) Search(ctx context.Context, sc queries.SearchCriteria) ([]T, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]T, error) {
		return s.core.Search(ctx, sc)
	})
}

func (s *Service[T]) Get(ctx context.Context, id any) (*T, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*T, error) {
		return s.core.Get(ctx, id)
	})
}
