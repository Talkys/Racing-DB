package articletext

import (
	"backend/internal/article"
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// Texts of hidden (draft or future-dated) articles are hidden too.
var ArticleTextSearchConfig = queries.EntitySearchConfig{
	BaseCondition: "article_id IN (SELECT id FROM articles WHERE " + article.VisibleCondition + ")",
	PageSize:      25,
	DefaultOrder:  "language ASC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":         queries.UUIDFilter("id"),
		"article_id": queries.UUIDFilter("article_id"),
		"language":   queries.TextFilter("language"),
		"title":      queries.TextFilter("title"),
	},
	AllowedSorts: map[string]string{
		"title":    "title",
		"language": "language",
	},
}

var articleTextColumns = []string{"id", "article_id", "language", "title", "summary", "body_md", "created_at", "updated_at"}

func scanArticleText(row database.Scanner) (ArticleText, error) {
	var x ArticleText
	err := row.Scan(&x.ID, &x.ArticleID, &x.Language, &x.Title, &x.Summary, &x.BodyMd, &x.CreatedAt, &x.UpdatedAt)
	return x, err
}

type ArticleTextRepository struct {
	cfg queries.EntitySearchConfig
}

func NewArticleTextRepository() *ArticleTextRepository {
	return &ArticleTextRepository{cfg: ArticleTextSearchConfig}
}

func (r *ArticleTextRepository) SearchArticleTexts(ctx context.Context, sc queries.SearchCriteria) ([]ArticleText, error) {
	rows, err := database.PerformSearch(ctx, sc, articleTextColumns, "article_texts", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []ArticleText{}

	for rows.Next() {
		x, err := scanArticleText(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *ArticleTextRepository) GetArticleText(ctx context.Context, id uuid.UUID) (*ArticleText, error) {
	row, err := database.PerformGet(ctx, id, articleTextColumns, "article_texts", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanArticleText(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("article text '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
