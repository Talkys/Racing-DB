package core

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/utils/queries"
	"context"
)

type CategoryCore struct {
	repo *repositories.CategoryRepository
}

func NewCategoryCore() *CategoryCore {
	return &CategoryCore{
		repo: repositories.NewCategoryRepository(),
	}
}

func (c *CategoryCore) SearchCategories(ctx context.Context, sc queries.SearchCriteria) ([]models.Category, error) {
	return c.repo.SearchCategories(ctx, sc)
}

func (c *CategoryCore) ListCategories(ctx context.Context) ([]models.Category, error) {
	return c.repo.ListActiveCategories(ctx)
}
