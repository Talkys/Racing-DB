package category

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"
)

var CategorySearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":          {DBColumn: "id", Kind: queries.KindText, AllowedOps: queries.IDOps},
		"name":        queries.TextFilter("name"),
		"description": queries.TextFilter("description"),
		"is_active":   queries.BoolFilter("is_active"),
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

var categoryColumns = []string{"id", "name", "description", "is_active"}

func scanCategory(row database.Scanner) (Category, error) {
	var x Category
	err := row.Scan(&x.ID, &x.Name, &x.Description, &x.Is_active)
	return x, err
}

type CategoryRepository struct {
	cfg queries.EntitySearchConfig
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{cfg: CategorySearchConfig}
}

func (r *CategoryRepository) SearchCategories(ctx context.Context, sc queries.SearchCriteria) ([]Category, error) {
	rows, err := database.PerformSearch(ctx, sc, categoryColumns, "categories", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Category{}

	for rows.Next() {
		x, err := scanCategory(rows)
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

func (r *CategoryRepository) GetCategory(ctx context.Context, id string) (*Category, error) {
	row, err := database.PerformGet(ctx, id, categoryColumns, "categories", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanCategory(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("category '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
