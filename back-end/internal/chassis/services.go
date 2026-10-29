package chassis

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type ChassisService struct {
	db   *sql.DB
	core *ChassisCore
}

func NewChassisService(db *sql.DB) *ChassisService {
	return &ChassisService{
		db:   db,
		core: NewChassisCore(),
	}
}

func (s *ChassisService) SearchChassis(ctx context.Context, sc queries.SearchCriteria) ([]Chassis, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Chassis, error) {
		return s.core.SearchChassis(ctx, sc)
	})
}

func (s *ChassisService) GetChassis(ctx context.Context, id uuid.UUID) (*Chassis, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Chassis, error) {
		return s.core.GetChassis(ctx, id)
	})
}
