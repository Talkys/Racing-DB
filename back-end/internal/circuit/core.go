package circuit

import (
	"backend/utils/queries"
	"context"
)

type CircuitCore struct {
	repo *CircuitRepository
}

func NewCircuitCore() *CircuitCore {
	return &CircuitCore{
		repo: NewCircuitRepository(),
	}
}

func (c *CircuitCore) SearchCircuits(ctx context.Context, sc queries.SearchCriteria) ([]Circuit, error) {
	return c.repo.SearchCircuits(ctx, sc)
}
