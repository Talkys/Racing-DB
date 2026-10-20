package chassis

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type ChassisService struct {
	db   *sql.DB
	core *ChassisCore
}

func NewChassisService(db *sql.DB) *ChassisService {
	return &ChassisService{
		db:   db,
		core: NewChassisCore(),
	}
}

func (s *ChassisService) SearchChassiss(ctx context.Context, sc queries.SearchCriteria) ([]Chassis, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchChassiss(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
