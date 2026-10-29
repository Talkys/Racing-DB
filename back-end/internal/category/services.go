package category

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type CategoryService struct {
	db   *sql.DB
	core *CategoryCore
}

func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{
		db:   db,
		core: NewCategoryCore(),
	}
}

func (s *CategoryService) SearchCategories(ctx context.Context, sc queries.SearchCriteria) ([]Category, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Category, error) {
		return s.core.SearchCategories(ctx, sc)
	})
}

func (s *CategoryService) GetCategory(ctx context.Context, id string) (*Category, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Category, error) {
		return s.core.GetCategory(ctx, id)
	})
}
