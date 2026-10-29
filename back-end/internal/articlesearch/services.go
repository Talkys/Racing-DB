package articlesearch

import (
	"backend/modules/database"
	"context"
	"database/sql"
)

type ArticleSearchService struct {
	db   *sql.DB
	core *ArticleSearchCore
}

func NewArticleSearchService(db *sql.DB) *ArticleSearchService {
	return &ArticleSearchService{
		db:   db,
		core: NewArticleSearchCore(),
	}
}

func (s *ArticleSearchService) SearchArticles(ctx context.Context, sc SearchCriteria) ([]ArticleSearchResult, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]ArticleSearchResult, error) {
		return s.core.SearchArticles(ctx, sc)
	})
}
