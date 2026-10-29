package articletext

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
	"database/sql"

	"github.com/google/uuid"
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
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) ([]ArticleText, error) {
		return s.core.SearchArticleTexts(ctx, sc)
	})
}

func (s *ArticleTextService) GetArticleText(ctx context.Context, id uuid.UUID) (*ArticleText, error) {
	return database.WithReadTx(ctx, s.db, func(ctx context.Context) (*ArticleText, error) {
		return s.core.GetArticleText(ctx, id)
	})
}
