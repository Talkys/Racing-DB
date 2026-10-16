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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	var categories []Category

	if categories, err = s.core.SearchCategories(ctxTx, sc); err != nil {
		return nil, err
	}
	tx.Commit()
	return categories, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	var categories []Category

	if categories, err = s.core.ListCategories(ctxTx); err != nil {
		return nil, err
	}
	tx.Commit()
	return categories, nil
}
