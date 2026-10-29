package driver

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type DriverService struct {
	db   *sql.DB
	core *DriverCore
}

func NewDriverService(db *sql.DB) *DriverService {
	return &DriverService{
		db:   db,
		core: NewDriverCore(),
	}
}

func (s *DriverService) SearchDrivers(ctx context.Context, sc queries.SearchCriteria) ([]Driver, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Driver, error) {
		return s.core.SearchDrivers(ctx, sc)
	})
}

func (s *DriverService) GetDriver(ctx context.Context, id uuid.UUID) (*Driver, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Driver, error) {
		return s.core.GetDriver(ctx, id)
	})
}
