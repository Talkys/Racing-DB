package driver

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var DriverSearchConfig = queries.EntitySearchConfig{
	PageSize:     25,
	DefaultOrder: "last_name ASC, first_name ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":          queries.UUIDFilter("id"),
		"first_name":  queries.TextFilter("first_name"),
		"last_name":   queries.TextFilter("last_name"),
		"code":        queries.TextFilter("code"),
		"nationality": queries.TextFilter("nationality"),
	},
	AllowedSorts: map[string]string{
		"last_name": "last_name",
	},
}

var driverColumns = []string{"id", "first_name", "last_name", "code", "nationality", "created_at"}

func scanDriver(row database.Scanner) (Driver, error) {
	var x Driver
	err := row.Scan(&x.ID, &x.FirstName, &x.LastName, &x.Code, &x.Nationality, &x.CreatedAt)
	return x, err
}

type DriverRepository struct {
	cfg queries.EntitySearchConfig
}

func NewDriverRepository() *DriverRepository {
	return &DriverRepository{cfg: DriverSearchConfig}
}

func (r *DriverRepository) SearchDrivers(ctx context.Context, sc queries.SearchCriteria) ([]Driver, error) {
	rows, err := database.PerformSearch(ctx, sc, driverColumns, "drivers", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Driver{}

	for rows.Next() {
		x, err := scanDriver(rows)
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

func (r *DriverRepository) GetDriver(ctx context.Context, id uuid.UUID) (*Driver, error) {
	row, err := database.PerformGet(ctx, id, driverColumns, "drivers", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanDriver(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("driver '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
