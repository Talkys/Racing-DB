package article

import (
	"backend/modules/database"
	"backend/utils/apperror"
	"backend/utils/queries"
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

// VisibleCondition hides drafts (no published_at) and future-dated articles everywhere in the API.
const VisibleCondition = "published_at IS NOT NULL AND published_at <= now()"

var ArticleSearchConfig = queries.EntitySearchConfig{
	BaseCondition: VisibleCondition,
	PageSize:      25,
	DefaultOrder:  "published_at DESC, id ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"id":          queries.UUIDFilter("id"),
		"entity_type": queries.TextFilter("entity_type"),
		"entity_id":   queries.TextFilter("entity_id"),
		"kind":        queries.TextFilter("kind"),
	},
	AllowedSorts: map[string]string{
		"published_at": "published_at",
		"created_at":   "created_at",
	},
}

var articleColumns = []string{"id", "entity_type", "entity_id", "kind", "published_at", "created_at"}

func scanArticle(row database.Scanner) (Article, error) {
	var x Article
	err := row.Scan(&x.ID, &x.EntityType, &x.EntityID, &x.Kind, &x.PublishedAt, &x.CreatedAt)
	return x, err
}

type ArticleRepository struct {
	cfg queries.EntitySearchConfig
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{cfg: ArticleSearchConfig}
}

func (r *ArticleRepository) SearchArticles(ctx context.Context, sc queries.SearchCriteria) ([]Article, error) {
	rows, err := database.PerformSearch(ctx, sc, articleColumns, "articles", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []Article{}

	for rows.Next() {
		x, err := scanArticle(rows)
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

func (r *ArticleRepository) GetArticle(ctx context.Context, id uuid.UUID) (*Article, error) {
	row, err := database.PerformGet(ctx, id, articleColumns, "articles", r.cfg)
	if err != nil {
		return nil, err
	}

	x, err := scanArticle(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.NotFound("article '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &x, nil
}
