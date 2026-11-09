package race

import (
	"encoding/json"
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Race struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID string          `json:"category_id"`
	CircuitID  uuid.UUID       `json:"circuit_id"`
	WinnerID   *uuid.UUID      `json:"winner_id"`
	Name       string          `json:"name"`
	Season     int             `json:"season"`
	Round      int             `json:"round"`
	RaceDate   time.Time       `json:"race_date"`
	Payload    json.RawMessage `json:"payload" swaggertype:"object"` // as stored, even when invalid
	CreatedAt  time.Time       `json:"created_at"`

	// PayloadValid is false when Payload breaks its category's contract (spec 3.3); PayloadErrors then says how.
	PayloadValid  bool     `json:"payload_valid"`
	PayloadErrors []string `json:"payload_errors,omitempty"`
}

var Resource = crud.Resource[Race]{
	Path:    "/races",
	Table:   "races",
	Name:    "race",
	ID:      crud.IDUUID,
	Columns: []string{"id", "category_id", "circuit_id", "winner_id", "name", "season", "round", "race_date", "payload", "created_at"},
	Scan: func(row database.Scanner) (Race, error) {
		var x Race
		err := row.Scan(&x.ID, &x.CategoryID, &x.CircuitID, &x.WinnerID, &x.Name, &x.Season, &x.Round, &x.RaceDate, &x.Payload, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "race_date DESC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":          queries.UUIDFilter("id"),
			"category_id": queries.TextFilter("category_id"),
			"circuit_id":  queries.UUIDFilter("circuit_id"),
			"winner_id":   queries.UUIDFilter("winner_id"),
			"name":        queries.TextFilter("name"),
			"season":      queries.IntFilter("season"),
			"round":       queries.IntFilter("round"),
		},
		AllowedSorts: map[string]string{
			"race_date":   "race_date",
			"name":        "name",
			"category_id": "category_id",
			"season":      "season",
			"round":       "round",
		},
	},
	Check: checkPayload,
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search races
// @Description Retrieve a page of races filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists. A payload that breaks its category's contract is returned as stored with payload_valid false.
// @Tags races
// @ID search-races
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query []string false "Filter expression column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |). Repeat it to AND several filters." collectionFormat(multi)
// @Param sort query []string false "Comma-separated columns, each optionally followed by asc or desc (e.g. race_date,desc), in priority order; may be repeated. Sortable: race_date, name, category_id, season, round" collectionFormat(multi)
// @Success 200 {array} Race "List of matching races"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /races [get]
func _() {}

// @Summary Get a race
// @Tags races
// @ID get-races
// @Produce json
// @Param id path string true "Race UUID"
// @Success 200 {object} Race
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /races/{id} [get]
func _() {}
