package driver

import (
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Driver struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Code        string    `json:"code"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
}

var Resource = crud.Resource[Driver]{
	Path:    "/drivers",
	Table:   "drivers",
	Name:    "driver",
	ID:      crud.IDUUID,
	Columns: []string{"id", "first_name", "last_name", "code", "nationality", "created_at"},
	Scan: func(row database.Scanner) (Driver, error) {
		var x Driver
		err := row.Scan(&x.ID, &x.FirstName, &x.LastName, &x.Code, &x.Nationality, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "last_name ASC, first_name ASC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":          queries.UUIDFilter("id"),
			"first_name":  queries.TextFilter("first_name"),
			"last_name":   queries.TextFilter("last_name"),
			"code":        queries.TextFilter("code"),
			"nationality": queries.TextFilter("nationality"),
		},
		AllowedSorts: map[string]string{
			"last_name": "last_name",
		},
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search drivers
// @Description Retrieve a page of drivers filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags drivers
// @ID search-drivers
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. last_name,asc)"
// @Success 200 {array} Driver "List of matching drivers"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /drivers [get]
func _() {}

// @Summary Get a driver
// @Tags drivers
// @ID get-drivers
// @Produce json
// @Param id path string true "Driver UUID"
// @Success 200 {object} Driver
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /drivers/{id} [get]
func _() {}
