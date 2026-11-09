package car

import (
	"encoding/json"
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Car struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID string          `json:"category_id"`
	TeamID     uuid.UUID       `json:"team_id"`
	EngineID   uuid.UUID       `json:"engine_id"`
	ChassisID  uuid.UUID       `json:"chassis_id"`
	Name       string          `json:"name"`
	Season     int             `json:"season"`
	Parts      json.RawMessage `json:"parts" swaggertype:"object"` // as stored, even when invalid
	CreatedAt  time.Time       `json:"created_at"`

	// PartsValid is false when Parts breaks the contract (spec 3.4); PartsErrors then says how.
	PartsValid  bool     `json:"parts_valid"`
	PartsErrors []string `json:"parts_errors,omitempty"`
}

var Resource = crud.Resource[Car]{
	Path:    "/cars",
	Table:   "cars",
	Name:    "car",
	ID:      crud.IDUUID,
	Columns: []string{"id", "category_id", "team_id", "engine_id", "chassis_id", "name", "season", "parts", "created_at"},
	Scan: func(row database.Scanner) (Car, error) {
		var x Car
		err := row.Scan(&x.ID, &x.CategoryID, &x.TeamID, &x.EngineID, &x.ChassisID, &x.Name, &x.Season, &x.Parts, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "season DESC, name ASC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":          queries.UUIDFilter("id"),
			"category_id": queries.TextFilter("category_id"),
			"team_id":     queries.UUIDFilter("team_id"),
			"engine_id":   queries.UUIDFilter("engine_id"),
			"chassis_id":  queries.UUIDFilter("chassis_id"),
			"name":        queries.TextFilter("name"),
			"season":      queries.IntFilter("season"),
		},
		AllowedSorts: map[string]string{
			"name":        "name",
			"season":      "season",
			"category_id": "category_id",
		},
	},
	Check: checkParts,
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search cars
// @Description Retrieve a page of cars filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists. Parts that break the contract are returned as stored with parts_valid false.
// @Tags cars
// @ID search-cars
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query []string false "Filter expression column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |). Repeat it to AND several filters." collectionFormat(multi)
// @Param sort query []string false "Comma-separated columns, each optionally followed by asc or desc (e.g. name,desc), in priority order; may be repeated. Sortable: name, season, category_id" collectionFormat(multi)
// @Success 200 {array} Car "List of matching cars"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /cars [get]
func _() {}

// @Summary Get a car
// @Tags cars
// @ID get-cars
// @Produce json
// @Param id path string true "Car UUID"
// @Success 200 {object} Car
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /cars/{id} [get]
func _() {}
