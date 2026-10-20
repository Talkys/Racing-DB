package team

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var TeamSearchConfig = queries.EntitySearchConfig{
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
		"short_name": {
			DBColumn: "short_name",
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
		"founded_year": {
			DBColumn: "founded_year",
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

type TeamRepository struct {
	cfg queries.EntitySearchConfig
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{cfg: TeamSearchConfig}
}

func (r *TeamRepository) SearchTeams(ctx context.Context, sc queries.SearchCriteria) ([]Team, error) {
	cols := []string{"id", "name", "short_name", "country", "founded_year", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "teams", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Team

	for rows.Next() {
		var x Team
		if err := rows.Scan(&x.ID, &x.Name, &x.ShortName, &x.Country, &x.FoundedYear, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
