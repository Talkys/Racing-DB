package articletext

import (
	"backend/utils/queries"
	"context"
)

type ArticleTextCore struct {
	repo *ArticleTextRepository
}

func NewArticleTextCore() *ArticleTextCore {
	return &ArticleTextCore{
		repo: NewArticleTextRepository(),
	}
}

func (c *ArticleTextCore) SearchArticleTexts(ctx context.Context, sc queries.SearchCriteria) ([]ArticleText, error) {
	return c.repo.SearchArticleTexts(ctx, sc)
}
