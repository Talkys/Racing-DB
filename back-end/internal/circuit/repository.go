package circuit

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var CircuitSearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":        queries.UUIDFilter("id"),
		"name":      queries.TextFilter("name"),
		"location":  queries.TextFilter("location"),
		"country":   queries.TextFilter("country"),
		"length_km": queries.DecimalFilter("length_km"),
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

var circuitColumns = []string{"id", "name", "location", "country", "length_km", "created_at"}

func scanCircuit(row database.Scanner) (Circuit, error) {
	var x Circuit
	err := row.Scan(&x.ID, &x.Name, &x.Location, &x.Country, &x.LengthKm, &x.CreatedAt)
	return x, err
}

type CircuitRepository struct {
	cfg queries.EntitySearchConfig
}

func NewCircuitRepository() *CircuitRepository {
	return &CircuitRepository{cfg: CircuitSearchConfig}
}

func (r *CircuitRepository) SearchCircuits(ctx context.Context, sc queries.SearchCriteria) ([]Circuit, error) {
	rows, err := database.PerformSearch(ctx, sc, circuitColumns, "circuits", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Circuit{}

	for rows.Next() {
		x, err := scanCircuit(rows)
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

func (r *CircuitRepository) GetCircuit(ctx context.Context, id uuid.UUID) (*Circuit, error) {
	row, err := database.PerformGet(ctx, id, circuitColumns, "circuits", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanCircuit(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("circuit '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
