package services

import (
	"backend/internal/core"
	"backend/internal/database"
	"backend/internal/models"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type CategoryService struct {
	db   *sql.DB
	core *core.CategoryCore
}

func NewCategoryService(db *sql.DB) *CategoryService {
	return &CategoryService{
		db:   db,
		core: core.NewCategoryCore(),
	}
}

func (s *CategoryService) SearchCategories(ctx context.Context, sc queries.SearchCriteria) ([]models.Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	var categories []models.Category

	if categories, err = s.core.SearchCategories(ctxTx, sc); err != nil {
		return nil, err
	}
	tx.Commit()
	return categories, nil
}

func (s *CategoryService) ListCategories(ctx context.Context) ([]models.Category, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	var categories []models.Category

	if categories, err = s.core.ListCategories(ctxTx); err != nil {
		return nil, err
	}
	tx.Commit()
	return categories, nil
}
