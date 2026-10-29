package driver

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
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

func (c *DriverCore) GetDriver(ctx context.Context, id uuid.UUID) (*Driver, error) {
	return c.repo.GetDriver(ctx, id)
}
