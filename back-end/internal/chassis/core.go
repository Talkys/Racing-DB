package chassis

import (
	"backend/utils/queries"
	"context"
)

type ChassisCore struct {
	repo *ChassisRepository
}

func NewChassisCore() *ChassisCore {
	return &ChassisCore{
		repo: NewChassisRepository(),
	}
}

func (c *ChassisCore) SearchChassiss(ctx context.Context, sc queries.SearchCriteria) ([]Chassis, error) {
	return c.repo.SearchChassiss(ctx, sc)
}
