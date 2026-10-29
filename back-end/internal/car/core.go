package car

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
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

func (c *CarCore) GetCar(ctx context.Context, id uuid.UUID) (*Car, error) {
	return c.repo.GetCar(ctx, id)
}
