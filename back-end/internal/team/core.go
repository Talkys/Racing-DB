package team

import (
	"backend/utils/queries"
	"context"
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
