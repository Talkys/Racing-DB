package article

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]Article, error) {
		return s.core.SearchArticles(ctx, sc)
	})
}

func (s *ArticleService) GetArticle(ctx context.Context, id uuid.UUID) (*Article, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*Article, error) {
		return s.core.GetArticle(ctx, id)
	})
}
