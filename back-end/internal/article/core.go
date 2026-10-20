package article

import (
	"backend/utils/queries"
	"context"
)

type ArticleCore struct {
	repo *ArticleRepository
}

func NewArticleCore() *ArticleCore {
	return &ArticleCore{
		repo: NewArticleRepository(),
	}
}

func (c *ArticleCore) SearchArticles(ctx context.Context, sc queries.SearchCriteria) ([]Article, error) {
	return c.repo.SearchArticles(ctx, sc)
}
