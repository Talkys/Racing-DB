package category

import (
	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"
)

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

var Resource = crud.Resource[Category]{
	Path:    "/categories",
	Table:   "categories",
	Name:    "category",
	ID:      crud.IDSlug,
	Columns: []string{"id", "name", "description", "is_active"},
	Scan: func(row database.Scanner) (Category, error) {
		var x Category
		err := row.Scan(&x.ID, &x.Name, &x.Description, &x.IsActive)
		return x, err
	},
	Search: queries.EntitySearchConfig{
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
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search categories
// @Description Retrieve a page of categories filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags categories
// @ID search-categories
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Category "List of matching categories"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /categories [get]
func _() {}

// @Summary Get a category
// @Tags categories
// @ID get-categories
// @Produce json
// @Param id path string true "Category slug, e.g. f1 (case-insensitive)"
// @Success 200 {object} Category
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /categories/{id} [get]
func _() {}
