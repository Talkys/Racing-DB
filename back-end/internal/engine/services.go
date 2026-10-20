package engine

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchEngines(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
