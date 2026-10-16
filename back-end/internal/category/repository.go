package category

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
	"fmt"
)

var CategorySearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"name": {
			DBColumn: "name",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "!=",
				constants.Contains: "ilike",
			},
		},
		"description": {
			DBColumn: "description",
			AllowedOps: map[string]string{
				constants.Contains: "ilike",
			},
		},
		"is_active": {
			DBColumn: "is_active",
			AllowedOps: map[string]string{
				constants.Equal: "=",
			},
		},
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

type CategoryRepository struct {
	cfg queries.EntitySearchConfig
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{cfg: CategorySearchConfig}
}

func (r *CategoryRepository) SearchCategories(ctx context.Context, sc queries.SearchCriteria) ([]Category, error) {
	exec := database.ExtractTx(ctx)
	if exec == nil {
		return nil, fmt.Errorf("No transaction provided")
	}

	cols := []string{"id", "name", "description", "is_active"}
	query, args, err := queries.BuildSearchQuery(cols, "categories", r.cfg, sc)

	if err != nil {
		return nil, err
	}

	rows, err := exec.QueryContext(ctx, query, args[0])
	if err != nil {
		return nil, err
	}

	var categories []Category

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Is_active); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) ListActiveCategories(ctx context.Context) ([]Category, error) {
	exec := database.ExtractTx(ctx)
	if exec == nil {
		return nil, fmt.Errorf("No transaction provided")
	}

	query := "select id, name, description, is_active from categories where is_active = true"
	rows, err := exec.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var categories []Category

	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.Is_active); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
