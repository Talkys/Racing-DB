package driver

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var DriverSearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"first_name": {
			DBColumn: "first_name",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"last_name": {
			DBColumn: "last_name",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"code": {
			DBColumn: "code",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"nationality": {
			DBColumn: "nationality",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
	},
	AllowedSorts: map[string]string{
		"last_name": "last_name",
	},
}

type DriverRepository struct {
	cfg queries.EntitySearchConfig
}

func NewDriverRepository() *DriverRepository {
	return &DriverRepository{cfg: DriverSearchConfig}
}

func (r *DriverRepository) SearchDrivers(ctx context.Context, sc queries.SearchCriteria) ([]Driver, error) {
	cols := []string{"id", "first_name", "last_name", "code", "nationality", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "drivers", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Driver

	for rows.Next() {
		var x Driver
		if err := rows.Scan(&x.ID, &x.FirstName, &x.LastName, &x.Code, &x.Nationality, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
