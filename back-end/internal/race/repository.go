package race

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var RaceSearchConfig = queries.EntitySearchConfig{
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
		"round": {
			DBColumn: "round",
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
		"race_date": "race_date",
	},
}

type RaceRepository struct {
	cfg queries.EntitySearchConfig
}

func NewRaceRepository() *RaceRepository {
	return &RaceRepository{cfg: RaceSearchConfig}
}

func (r *RaceRepository) SearchRaces(ctx context.Context, sc queries.SearchCriteria) ([]Race, error) {
	cols := []string{"id", "category_id", "circuit_id", "winner_id", "name", "season", "round", "race_date", "payload", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "races", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Race

	for rows.Next() {
		var x Race
		if err := rows.Scan(&x.ID, &x.CategoryID, &x.CircuitID, &x.WinnerID, &x.Name, &x.Season, &x.Round, &x.RaceDate, &x.Payload, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
