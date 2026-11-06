package engine

import (
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Engine struct {
	ID             uuid.UUID `json:"id"`
	Manufacturer   string    `json:"manufacturer"`
	Name           string    `json:"name"`
	Configuration  string    `json:"configuration"`
	DisplacementCC int       `json:"displacement_cc"`
	Cylinders      int16     `json:"cylinders"`
	Aspiration     string    `json:"aspiration"`
	PowerKW        *float64  `json:"power_kw"`
	CreatedAt      time.Time `json:"created_at"`
}

var Resource = crud.Resource[Engine]{
	Path:    "/engines",
	Table:   "engines",
	Name:    "engine",
	ID:      crud.IDUUID,
	Columns: []string{"id", "manufacturer", "name", "configuration", "displacement_cc", "cylinders", "aspiration", "power_kw", "created_at"},
	Scan: func(row database.Scanner) (Engine, error) {
		var x Engine
		err := row.Scan(&x.ID, &x.Manufacturer, &x.Name, &x.Configuration, &x.DisplacementCC, &x.Cylinders, &x.Aspiration, &x.PowerKW, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "manufacturer ASC, name ASC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":              queries.UUIDFilter("id"),
			"manufacturer":    queries.TextFilter("manufacturer"),
			"name":            queries.TextFilter("name"),
			"configuration":   queries.TextFilter("configuration"),
			"displacement_cc": queries.IntFilter("displacement_cc"),
			"cylinders":       queries.IntFilter("cylinders"),
			"aspiration":      queries.TextFilter("aspiration"),
			"power_kw":        queries.DecimalFilter("power_kw"),
		},
		AllowedSorts: map[string]string{
			"name": "name",
		},
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search engines
// @Description Retrieve a page of engines filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags engines
// @ID search-engines
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Engine "List of matching engines"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /engines [get]
func _() {}

// @Summary Get an engine
// @Tags engines
// @ID get-engines
// @Produce json
// @Param id path string true "Engine UUID"
// @Success 200 {object} Engine
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /engines/{id} [get]
func _() {}
