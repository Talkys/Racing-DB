package race

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchRaces(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
