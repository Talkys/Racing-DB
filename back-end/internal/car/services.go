package car

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type CarService struct {
	db   *sql.DB
	core *CarCore
}

func NewCarService(db *sql.DB) *CarService {
	return &CarService{
		db:   db,
		core: NewCarCore(),
	}
}

func (s *CarService) SearchCars(ctx context.Context, sc queries.SearchCriteria) ([]Car, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Car, error) {
		return s.core.SearchCars(ctx, sc)
	})
}

func (s *CarService) GetCar(ctx context.Context, id uuid.UUID) (*Car, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Car, error) {
		return s.core.GetCar(ctx, id)
	})
}
