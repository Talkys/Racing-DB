package circuit

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var CircuitSearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"name": {
			DBColumn: "name",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"location": {
			DBColumn: "location",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"country": {
			DBColumn: "country",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"length_km": {
			DBColumn: "length_km",
			AllowedOps: map[string]string{
				constants.Equal:        "=",
				constants.Greater:      ">",
				constants.GreaterEqual: ">=",
				constants.Less:         "<",
				constants.LessEqual:    "<=",
			},
		},
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

type CircuitRepository struct {
	cfg queries.EntitySearchConfig
}

func NewCircuitRepository() *CircuitRepository {
	return &CircuitRepository{cfg: CircuitSearchConfig}
}

func (r *CircuitRepository) SearchCircuits(ctx context.Context, sc queries.SearchCriteria) ([]Circuit, error) {
	cols := []string{"id", "name", "location", "country", "length_km", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "circuits", r.cfg)
	if err != nil {
		return nil, err
	}

	var categories []Circuit

	for rows.Next() {
		var c Circuit
		if err := rows.Scan(&c.ID, &c.Name, &c.Location, &c.Country, &c.LengthKm, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
