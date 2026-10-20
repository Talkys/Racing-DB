package car

import (
	"backend/utils/queries"
	"context"
)

type CarCore struct {
	repo *CarRepository
}

func NewCarCore() *CarCore {
	return &CarCore{
		repo: NewCarRepository(),
	}
}

func (c *CarCore) SearchCars(ctx context.Context, sc queries.SearchCriteria) ([]Car, error) {
	return c.repo.SearchCars(ctx, sc)
}
