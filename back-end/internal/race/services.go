package race

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type RaceService struct {
	db   *sql.DB
	core *RaceCore
}

func NewRaceService(db *sql.DB) *RaceService {
	return &RaceService{
		db:   db,
		core: NewRaceCore(),
	}
}

func (s *RaceService) SearchRaces(ctx context.Context, sc queries.SearchCriteria) ([]Race, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Race, error) {
		return s.core.SearchRaces(ctx, sc)
	})
}

func (s *RaceService) GetRace(ctx context.Context, id uuid.UUID) (*Race, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Race, error) {
		return s.core.GetRace(ctx, id)
	})
}
