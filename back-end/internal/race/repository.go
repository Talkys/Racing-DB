package race

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var RaceSearchConfig = queries.EntitySearchConfig{
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
		"race_date": "race_date",
	},
}

var raceColumns = []string{"id", "category_id", "circuit_id", "winner_id", "name", "season", "round", "race_date", "payload", "created_at"}

func scanRace(row database.Scanner) (Race, error) {
	var x Race
	err := row.Scan(&x.ID, &x.CategoryID, &x.CircuitID, &x.WinnerID, &x.Name, &x.Season, &x.Round, &x.RaceDate, &x.Payload, &x.CreatedAt)
	return x, err
}

type RaceRepository struct {
	cfg queries.EntitySearchConfig
}

func NewRaceRepository() *RaceRepository {
	return &RaceRepository{cfg: RaceSearchConfig}
}

func (r *RaceRepository) SearchRaces(ctx context.Context, sc queries.SearchCriteria) ([]Race, error) {
	rows, err := database.PerformSearch(ctx, sc, raceColumns, "races", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Race{}

	for rows.Next() {
		x, err := scanRace(rows)
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

func (r *RaceRepository) GetRace(ctx context.Context, id uuid.UUID) (*Race, error) {
	row, err := database.PerformGet(ctx, id, raceColumns, "races", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanRace(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("race '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
