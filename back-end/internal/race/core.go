package race

import (
	"backend/utils/queries"
	"context"
)

type RaceCore struct {
	repo *RaceRepository
}

func NewRaceCore() *RaceCore {
	return &RaceCore{
		repo: NewRaceRepository(),
	}
}

func (c *RaceCore) SearchRaces(ctx context.Context, sc queries.SearchCriteria) ([]Race, error) {
	return c.repo.SearchRaces(ctx, sc)
}
