package team

import (
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Team struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ShortName   string    `json:"short_name"`
	Country     string    `json:"country"`
	FoundedYear *int      `json:"founded_year"`
	CreatedAt   time.Time `json:"created_at"`
}

var Resource = crud.Resource[Team]{
	Path:    "/teams",
	Table:   "teams",
	Name:    "team",
	ID:      crud.IDUUID,
	Columns: []string{"id", "name", "short_name", "country", "founded_year", "created_at"},
	Scan: func(row database.Scanner) (Team, error) {
		var x Team
		err := row.Scan(&x.ID, &x.Name, &x.ShortName, &x.Country, &x.FoundedYear, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "name ASC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":           queries.UUIDFilter("id"),
			"name":         queries.TextFilter("name"),
			"short_name":   queries.TextFilter("short_name"),
			"country":      queries.TextFilter("country"),
			"founded_year": queries.IntFilter("founded_year"),
		},
		AllowedSorts: map[string]string{
			"name": "name",
		},
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search teams
// @Description Retrieve a page of teams filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags teams
// @ID search-teams
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Team "List of matching teams"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /teams [get]
func _() {}

// @Summary Get a team
// @Tags teams
// @ID get-teams
// @Produce json
// @Param id path string true "Team UUID"
// @Success 200 {object} Team
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /teams/{id} [get]
func _() {}
