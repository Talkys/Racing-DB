package article

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"
)

type ArticleService struct {
	db   *sql.DB
	core *ArticleCore
}

func NewArticleService(db *sql.DB) *ArticleService {
	return &ArticleService{
		db:   db,
		core: NewArticleCore(),
	}
}

func (s *ArticleService) SearchArticles(ctx context.Context, sc queries.SearchCriteria) ([]Article, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	ctxTx := database.InjectTx(ctx, tx)

	results, err := s.core.SearchArticles(ctxTx, sc)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return results, nil
}
