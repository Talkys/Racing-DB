package category

import (
	"backend/utils/queries"
	"context"
)

type CategoryCore struct {
	repo *CategoryRepository
}

func NewCategoryCore() *CategoryCore {
	return &CategoryCore{
		repo: NewCategoryRepository(),
	}
}

func (c *CategoryCore) SearchCategories(ctx context.Context, sc queries.SearchCriteria) ([]Category, error) {
	return c.repo.SearchCategories(ctx, sc)
}

func (c *CategoryCore) GetCategory(ctx context.Context, id string) (*Category, error) {
	return c.repo.GetCategory(ctx, id)
}
