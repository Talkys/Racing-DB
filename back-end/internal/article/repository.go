package article

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var ArticleSearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"entity_type": {
			DBColumn: "entity_type",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"entity_id": {
			DBColumn: "entity_id",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"kind": {
			DBColumn: "kind",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
	},
	AllowedSorts: map[string]string{
		"created_at": "created_at",
	},
}

type ArticleRepository struct {
	cfg queries.EntitySearchConfig
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{cfg: ArticleSearchConfig}
}

func (r *ArticleRepository) SearchArticles(ctx context.Context, sc queries.SearchCriteria) ([]Article, error) {
	cols := []string{"id", "entity_type", "entity_id", "kind", "published_at", "created_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "articles", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Article

	for rows.Next() {
		var x Article
		if err := rows.Scan(&x.ID, &x.EntityType, &x.EntityID, &x.Kind, &x.PublishedAt, &x.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
