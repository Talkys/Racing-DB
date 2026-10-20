package driver

import (
	"backend/utils/queries"
	"context"
)

type DriverCore struct {
	repo *DriverRepository
}

func NewDriverCore() *DriverCore {
	return &DriverCore{
		repo: NewDriverRepository(),
	}
}

func (c *DriverCore) SearchDrivers(ctx context.Context, sc queries.SearchCriteria) ([]Driver, error) {
	return c.repo.SearchDrivers(ctx, sc)
}
