package engine

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
)

type EngineCore struct {
	repo *EngineRepository
}

func NewEngineCore() *EngineCore {
	return &EngineCore{
		repo: NewEngineRepository(),
	}
}

func (c *EngineCore) SearchEngines(ctx context.Context, sc queries.SearchCriteria) ([]Engine, error) {
	return c.repo.SearchEngines(ctx, sc)
}

func (c *EngineCore) GetEngine(ctx context.Context, id uuid.UUID) (*Engine, error) {
	return c.repo.GetEngine(ctx, id)
}
