package car

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type CarService struct {
	db   *sql.DB
	core *CarCore
}

func NewCarService(db *sql.DB) *CarService {
	return &CarService{
		db:   db,
		core: NewCarCore(),
	}
}

func (s *CarService) SearchCars(ctx context.Context, sc queries.SearchCriteria) ([]Car, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchCars(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
