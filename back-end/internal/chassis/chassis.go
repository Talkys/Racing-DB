package chassis

import (
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Chassis struct {
	ID           uuid.UUID `json:"id"`
	Manufacturer string    `json:"manufacturer"`
	Name         string    `json:"name"`
	Material     string    `json:"material"`
	WheelbaseMM  *int      `json:"wheelbase_mm"`
	WeightKG     *float64  `json:"weight_kg"`
	CreatedAt    time.Time `json:"created_at"`
}

var Resource = crud.Resource[Chassis]{
	Path:    "/chassis",
	Table:   "chassis",
	Name:    "chassis",
	ID:      crud.IDUUID,
	Columns: []string{"id", "manufacturer", "name", "material", "wheelbase_mm", "weight_kg", "created_at"},
	Scan: func(row database.Scanner) (Chassis, error) {
		var x Chassis
		err := row.Scan(&x.ID, &x.Manufacturer, &x.Name, &x.Material, &x.WheelbaseMM, &x.WeightKG, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
		PageSize:     25,
		DefaultOrder: "manufacturer ASC, name ASC, id ASC",
		AllowedFilters: map[string]queries.AllowedFilter{
			"id":           queries.UUIDFilter("id"),
			"manufacturer": queries.TextFilter("manufacturer"),
			"name":         queries.TextFilter("name"),
			"material":     queries.TextFilter("material"),
			"wheelbase_mm": queries.IntFilter("wheelbase_mm"),
			"weight_kg":    queries.DecimalFilter("weight_kg"),
		},
		AllowedSorts: map[string]string{
			"name": "name",
		},
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search chassis
// @Description Retrieve a page of chassis filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags chassis
// @ID search-chassis
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Chassis "List of matching chassis"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /chassis [get]
func _() {}

// @Summary Get a chassis
// @Tags chassis
// @ID get-chassis
// @Produce json
// @Param id path string true "Chassis UUID"
// @Success 200 {object} Chassis
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /chassis/{id} [get]
func _() {}
