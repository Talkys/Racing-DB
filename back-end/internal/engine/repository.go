package engine

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var EngineSearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "manufacturer ASC, name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":              queries.UUIDFilter("id"),
		"manufacturer":    queries.TextFilter("manufacturer"),
		"name":            queries.TextFilter("name"),
		"configuration":   queries.TextFilter("configuration"),
		"displacement_cc": queries.IntFilter("displacement_cc"),
		"cylinders":       queries.IntFilter("cylinders"),
		"aspiration":      queries.TextFilter("aspiration"),
		"power_kw":        queries.DecimalFilter("power_kw"),
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

var engineColumns = []string{"id", "manufacturer", "name", "configuration", "displacement_cc", "cylinders", "aspiration", "power_kw", "created_at"}

func scanEngine(row database.Scanner) (Engine, error) {
	var x Engine
	err := row.Scan(&x.ID, &x.Manufacturer, &x.Name, &x.Configuration, &x.DisplacementCC, &x.Cylinders, &x.Aspiration, &x.PowerKW, &x.CreatedAt)
	return x, err
}

type EngineRepository struct {
	cfg queries.EntitySearchConfig
}

func NewEngineRepository() *EngineRepository {
	return &EngineRepository{cfg: EngineSearchConfig}
}

func (r *EngineRepository) SearchEngines(ctx context.Context, sc queries.SearchCriteria) ([]Engine, error) {
	rows, err := database.PerformSearch(ctx, sc, engineColumns, "engines", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Engine{}

	for rows.Next() {
		x, err := scanEngine(rows)
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

func (r *EngineRepository) GetEngine(ctx context.Context, id uuid.UUID) (*Engine, error) {
	row, err := database.PerformGet(ctx, id, engineColumns, "engines", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanEngine(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("engine '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
