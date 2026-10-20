package car

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var CarSearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"category_id": {
			DBColumn: "category_id",
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
		"season": {
			DBColumn: "season",
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

type CarRepository struct {
	cfg queries.EntitySearchConfig
}

func NewCarRepository() *CarRepository {
	return &CarRepository{cfg: CarSearchConfig}
}

func (r *CarRepository) SearchCars(ctx context.Context, sc queries.SearchCriteria) ([]Car, error) {
	cols := []string{"id", "category_id", "team_id", "engine_id", "chassis_id", "name", "season", "parts", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "cars", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Car

	for rows.Next() {
		var x Car
		if err := rows.Scan(&x.ID, &x.CategoryID, &x.TeamID, &x.EngineID, &x.ChassisID, &x.Name, &x.Season, &x.Parts, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
