package chassis

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var ChassisSearchConfig = queries.EntitySearchConfig{
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
		"material": {
			DBColumn: "material",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"wheelbase_mm": {
			DBColumn: "wheelbase_mm",
			AllowedOps: map[string]string{
				constants.Equal:        "=",
				constants.Greater:      ">",
				constants.GreaterEqual: ">=",
				constants.Less:         "<",
				constants.LessEqual:    "<=",
			},
		},
		"weight_kg": {
			DBColumn: "weight_kg",
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

type ChassisRepository struct {
	cfg queries.EntitySearchConfig
}

func NewChassisRepository() *ChassisRepository {
	return &ChassisRepository{cfg: ChassisSearchConfig}
}

func (r *ChassisRepository) SearchChassiss(ctx context.Context, sc queries.SearchCriteria) ([]Chassis, error) {
	cols := []string{"id", "manufacturer", "name", "material", "wheelbase_mm", "weight_kg", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "chassis", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Chassis

	for rows.Next() {
		var x Chassis
		if err := rows.Scan(&x.ID, &x.Manufacturer, &x.Name, &x.Material, &x.WheelbaseMM, &x.WeightKG, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
