package driver

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type DriverService struct {
	db   *sql.DB
	core *DriverCore
}

func NewDriverService(db *sql.DB) *DriverService {
	return &DriverService{
		db:   db,
		core: NewDriverCore(),
	}
}

func (s *DriverService) SearchDrivers(ctx context.Context, sc queries.SearchCriteria) ([]Driver, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchDrivers(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
