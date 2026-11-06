package articletext

import (
	"time"

	"backend/internal/article"
	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type ArticleText struct {
	ID        uuid.UUID `json:"id"`
	ArticleID uuid.UUID `json:"article_id"`
	Language  string    `json:"language"`
	Title     string    `json:"title"`
	Summary   *string   `json:"summary"`
	BodyMd    string    `json:"body_md"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var Resource = crud.Resource[ArticleText]{
	Path:    "/article-texts",
	Table:   "article_texts",
	Name:    "article text",
	ID:      crud.IDUUID,
	Columns: []string{"id", "article_id", "language", "title", "summary", "body_md", "created_at", "updated_at"},
	Scan: func(row database.Scanner) (ArticleText, error) {
		var x ArticleText
		err := row.Scan(&x.ID, &x.ArticleID, &x.Language, &x.Title, &x.Summary, &x.BodyMd, &x.CreatedAt, &x.UpdatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
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
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search article texts
// @Description Retrieve a page of article texts filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags article-texts
// @ID search-article-texts
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. title,asc)"
// @Success 200 {array} ArticleText "List of matching article texts"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /article-texts [get]
func _() {}

// @Summary Get an article text
// @Tags article-texts
// @ID get-article-texts
// @Produce json
// @Param id path string true "ArticleText UUID"
// @Success 200 {object} ArticleText
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /article-texts/{id} [get]
func _() {}
