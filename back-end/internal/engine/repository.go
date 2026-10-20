package engine

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var EngineSearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"manufacturer": {
			DBColumn: "manufacturer",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"name": {
			DBColumn: "name",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"configuration": {
			DBColumn: "configuration",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"displacement_cc": {
			DBColumn: "displacement_cc",
			AllowedOps: map[string]string{
				constants.Equal:        "=",
				constants.Greater:      ">",
				constants.GreaterEqual: ">=",
				constants.Less:         "<",
				constants.LessEqual:    "<=",
			},
		},
		"cylinders": {
			DBColumn: "cylinders",
			AllowedOps: map[string]string{
				constants.Equal:        "=",
				constants.Greater:      ">",
				constants.GreaterEqual: ">=",
				constants.Less:         "<",
				constants.LessEqual:    "<=",
			},
		},
		"aspiration": {
			DBColumn: "aspiration",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"power_kw": {
			DBColumn: "power_kw",
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

type EngineRepository struct {
	cfg queries.EntitySearchConfig
}

func NewEngineRepository() *EngineRepository {
	return &EngineRepository{cfg: EngineSearchConfig}
}

func (r *EngineRepository) SearchEngines(ctx context.Context, sc queries.SearchCriteria) ([]Engine, error) {
	cols := []string{"id", "manufacturer", "name", "configuration", "displacement_cc", "cylinders", "aspiration", "power_kw", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "engines", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Engine

	for rows.Next() {
		var x Engine
		if err := rows.Scan(&x.ID, &x.Manufacturer, &x.Name, &x.Configuration, &x.DisplacementCC, &x.Cylinders, &x.Aspiration, &x.PowerKW, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
