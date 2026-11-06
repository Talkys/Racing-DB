package circuit

import (
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Circuit struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Country   string    `json:"country"`
	LengthKm  float64   `json:"length_km"`
	CreatedAt time.Time `json:"created_at"`
}

var Resource = crud.Resource[Circuit]{
	Path:    "/circuits",
	Table:   "circuits",
	Name:    "circuit",
	ID:      crud.IDUUID,
	Columns: []string{"id", "name", "location", "country", "length_km", "created_at"},
	Scan: func(row database.Scanner) (Circuit, error) {
		var x Circuit
		err := row.Scan(&x.ID, &x.Name, &x.Location, &x.Country, &x.LengthKm, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "name ASC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":        queries.UUIDFilter("id"),
			"name":      queries.TextFilter("name"),
			"location":  queries.TextFilter("location"),
			"country":   queries.TextFilter("country"),
			"length_km": queries.DecimalFilter("length_km"),
		},
		AllowedSorts: map[string]string{
			"name": "name",
		},
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search circuits
// @Description Retrieve a page of circuits filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags circuits
// @ID search-circuits
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Circuit "List of matching circuits"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /circuits [get]
func _() {}

// @Summary Get a circuit
// @Tags circuits
// @ID get-circuits
// @Produce json
// @Param id path string true "Circuit UUID"
// @Success 200 {object} Circuit
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /circuits/{id} [get]
func _() {}
