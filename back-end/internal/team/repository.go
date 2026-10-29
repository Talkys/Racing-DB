package team

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var TeamSearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":           queries.UUIDFilter("id"),
		"name":         queries.TextFilter("name"),
		"short_name":   queries.TextFilter("short_name"),
		"country":      queries.TextFilter("country"),
		"founded_year": queries.IntFilter("founded_year"),
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

var teamColumns = []string{"id", "name", "short_name", "country", "founded_year", "created_at"}

func scanTeam(row database.Scanner) (Team, error) {
	var x Team
	err := row.Scan(&x.ID, &x.Name, &x.ShortName, &x.Country, &x.FoundedYear, &x.CreatedAt)
	return x, err
}

type TeamRepository struct {
	cfg queries.EntitySearchConfig
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{cfg: TeamSearchConfig}
}

func (r *TeamRepository) SearchTeams(ctx context.Context, sc queries.SearchCriteria) ([]Team, error) {
	rows, err := database.PerformSearch(ctx, sc, teamColumns, "teams", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Team{}

	for rows.Next() {
		x, err := scanTeam(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *TeamRepository) GetTeam(ctx context.Context, id uuid.UUID) (*Team, error) {
	row, err := database.PerformGet(ctx, id, teamColumns, "teams", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanTeam(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("team '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
