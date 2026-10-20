package team

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type TeamService struct {
	db   *sql.DB
	core *TeamCore
}

func NewTeamService(db *sql.DB) *TeamService {
	return &TeamService{
		db:   db,
		core: NewTeamCore(),
	}
}

func (s *TeamService) SearchTeams(ctx context.Context, sc queries.SearchCriteria) ([]Team, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchTeams(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
