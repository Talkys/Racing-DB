package circuit

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
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

func (c *CircuitCore) GetCircuit(ctx context.Context, id uuid.UUID) (*Circuit, error) {
	return c.repo.GetCircuit(ctx, id)
}
