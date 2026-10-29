package articletext

import (
	"backend/utils/queries"
	"context"

	"github.com/google/uuid"
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

func (c *ArticleTextCore) GetArticleText(ctx context.Context, id uuid.UUID) (*ArticleText, error) {
	return c.repo.GetArticleText(ctx, id)
}
