package engine

import (
	"backend/utils/queries"
	"context"
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
