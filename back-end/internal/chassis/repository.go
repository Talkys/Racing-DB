package chassis

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var ChassisSearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "manufacturer ASC, name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":           queries.UUIDFilter("id"),
		"manufacturer": queries.TextFilter("manufacturer"),
		"name":         queries.TextFilter("name"),
		"material":     queries.TextFilter("material"),
		"wheelbase_mm": queries.IntFilter("wheelbase_mm"),
		"weight_kg":    queries.DecimalFilter("weight_kg"),
	},
	AllowedSorts: map[string]string{
		"name": "name",
	},
}

var chassisColumns = []string{"id", "manufacturer", "name", "material", "wheelbase_mm", "weight_kg", "created_at"}

func scanChassis(row database.Scanner) (Chassis, error) {
	var x Chassis
	err := row.Scan(&x.ID, &x.Manufacturer, &x.Name, &x.Material, &x.WheelbaseMM, &x.WeightKG, &x.CreatedAt)
	return x, err
}

type ChassisRepository struct {
	cfg queries.EntitySearchConfig
}

func NewChassisRepository() *ChassisRepository {
	return &ChassisRepository{cfg: ChassisSearchConfig}
}

func (r *ChassisRepository) SearchChassis(ctx context.Context, sc queries.SearchCriteria) ([]Chassis, error) {
	rows, err := database.PerformSearch(ctx, sc, chassisColumns, "chassis", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Chassis{}

	for rows.Next() {
		x, err := scanChassis(rows)
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

func (r *ChassisRepository) GetChassis(ctx context.Context, id uuid.UUID) (*Chassis, error) {
	row, err := database.PerformGet(ctx, id, chassisColumns, "chassis", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanChassis(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("chassis '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
