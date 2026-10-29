package race

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
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

func (c *RaceCore) GetRace(ctx context.Context, id uuid.UUID) (*Race, error) {
	return c.repo.GetRace(ctx, id)
}
