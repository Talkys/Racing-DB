package articletext

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type ArticleTextService struct {
	db   *sql.DB
	core *ArticleTextCore
}

func NewArticleTextService(db *sql.DB) *ArticleTextService {
	return &ArticleTextService{
		db:   db,
		core: NewArticleTextCore(),
	}
}

func (s *ArticleTextService) SearchArticleTexts(ctx context.Context, sc queries.SearchCriteria) ([]ArticleText, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchArticleTexts(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
