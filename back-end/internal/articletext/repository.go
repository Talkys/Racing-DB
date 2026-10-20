package articletext

import (
	"backend/modules/database"
	"backend/utils/constants"
	"backend/utils/queries"
	"context"
)

var ArticleTextSearchConfig = queries.EntitySearchConfig{
	PageSize:       25,
	DefaultSortCol: "id",
	DefaultSortDir: "ASC",
	AllowedFilters: map[string]queries.AllowedFilter{
		"language": {
			DBColumn: "language",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
		"title": {
			DBColumn: "title",
			AllowedOps: map[string]string{
				constants.Equal:    "=",
				constants.NotEqual: "<>",
				constants.Contains: "ilike",
			},
		},
	},
	AllowedSorts: map[string]string{
		"title": "title",
	},
}

type ArticleTextRepository struct {
	cfg queries.EntitySearchConfig
}

func NewArticleTextRepository() *ArticleTextRepository {
	return &ArticleTextRepository{cfg: ArticleTextSearchConfig}
}

func (r *ArticleTextRepository) SearchArticleTexts(ctx context.Context, sc queries.SearchCriteria) ([]ArticleText, error) {
	cols := []string{"id", "article_id", "language", "title", "summary", "body_md", "created_at", "updated_at"}

	rows, err := database.PerformSearch(ctx, sc, cols, "article_texts", r.cfg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ArticleText

	for rows.Next() {
		var x ArticleText
		if err := rows.Scan(&x.ID, &x.ArticleID, &x.Language, &x.Title, &x.Summary, &x.BodyMd, &x.CreatedAt, &x.UpdatedAt); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
