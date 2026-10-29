package engine

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type EngineService struct {
	db   *sql.DB
	core *EngineCore
}

func NewEngineService(db *sql.DB) *EngineService {
	return &EngineService{
		db:   db,
		core: NewEngineCore(),
	}
}

func (s *EngineService) SearchEngines(ctx context.Context, sc queries.SearchCriteria) ([]Engine, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Engine, error) {
		return s.core.SearchEngines(ctx, sc)
	})
}

func (s *EngineService) GetEngine(ctx context.Context, id uuid.UUID) (*Engine, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Engine, error) {
		return s.core.GetEngine(ctx, id)
	})
}
