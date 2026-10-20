package circuit

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	var Circuits []Circuit

	if Circuits, err = s.core.SearchCircuits(ctxTx, sc); err != nil {
		return nil, err
	}
	tx.Commit()
	return Circuits, nil
}
