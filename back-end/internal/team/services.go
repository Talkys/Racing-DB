package team

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Team, error) {
		return s.core.SearchTeams(ctx, sc)
	})
}

func (s *TeamService) GetTeam(ctx context.Context, id uuid.UUID) (*Team, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Team, error) {
		return s.core.GetTeam(ctx, id)
	})
}
