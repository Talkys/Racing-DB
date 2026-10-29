package circuit

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type CircuitService struct {
	db   *sql.DB
	core *CircuitCore
}

func NewCircuitService(db *sql.DB) *CircuitService {
	return &CircuitService{
		db:   db,
		core: NewCircuitCore(),
	}
}

func (s *CircuitService) SearchCircuits(ctx context.Context, sc queries.SearchCriteria) ([]Circuit, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Circuit, error) {
		return s.core.SearchCircuits(ctx, sc)
	})
}

func (s *CircuitService) GetCircuit(ctx context.Context, id uuid.UUID) (*Circuit, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Circuit, error) {
		return s.core.GetCircuit(ctx, id)
	})
}
