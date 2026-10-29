package chassis

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
)

type ChassisCore struct {
	repo *ChassisRepository
}

func NewChassisCore() *ChassisCore {
	return &ChassisCore{
		repo: NewChassisRepository(),
	}
}

func (c *ChassisCore) SearchChassis(ctx context.Context, sc queries.SearchCriteria) ([]Chassis, error) {
	return c.repo.SearchChassis(ctx, sc)
}

func (c *ChassisCore) GetChassis(ctx context.Context, id uuid.UUID) (*Chassis, error) {
	return c.repo.GetChassis(ctx, id)
}
