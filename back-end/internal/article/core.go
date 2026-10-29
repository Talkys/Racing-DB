package article

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
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

func (c *ArticleCore) GetArticle(ctx context.Context, id uuid.UUID) (*Article, error) {
	return c.repo.GetArticle(ctx, id)
}
