package car

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var CarSearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "season DESC, name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":          queries.UUIDFilter("id"),
		"category_id": queries.TextFilter("category_id"),
		"team_id":     queries.UUIDFilter("team_id"),
		"engine_id":   queries.UUIDFilter("engine_id"),
		"chassis_id":  queries.UUIDFilter("chassis_id"),
		"name":        queries.TextFilter("name"),
		"season":      queries.IntFilter("season"),
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

var carColumns = []string{"id", "category_id", "team_id", "engine_id", "chassis_id", "name", "season", "parts", "created_at"}

func scanCar(row database.Scanner) (Car, error) {
	var x Car
	err := row.Scan(&x.ID, &x.CategoryID, &x.TeamID, &x.EngineID, &x.ChassisID, &x.Name, &x.Season, &x.Parts, &x.CreatedAt)
	return x, err
}

type CarRepository struct {
	cfg queries.EntitySearchConfig
}

func NewCarRepository() *CarRepository {
	return &CarRepository{cfg: CarSearchConfig}
}

func (r *CarRepository) SearchCars(ctx context.Context, sc queries.SearchCriteria) ([]Car, error) {
	rows, err := database.PerformSearch(ctx, sc, carColumns, "cars", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Car{}

	for rows.Next() {
		x, err := scanCar(rows)
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

func (r *CarRepository) GetCar(ctx context.Context, id uuid.UUID) (*Car, error) {
	row, err := database.PerformGet(ctx, id, carColumns, "cars", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanCar(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("car '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
