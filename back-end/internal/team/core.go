package team

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
)

type TeamCore struct {
	repo *TeamRepository
}

func NewTeamCore() *TeamCore {
	return &TeamCore{
		repo: NewTeamRepository(),
	}
}

func (c *TeamCore) SearchTeams(ctx context.Context, sc queries.SearchCriteria) ([]Team, error) {
	return c.repo.SearchTeams(ctx, sc)
}

func (c *TeamCore) GetTeam(ctx context.Context, id uuid.UUID) (*Team, error) {
	return c.repo.GetTeam(ctx, id)
}
