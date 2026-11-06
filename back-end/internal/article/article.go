package article

import (
	"time"

	"backend/internal/crud"
	"backend/modules/database"
	"backend/utils/queries"

	"github.com/google/uuid"
)

type Article struct {
	ID          uuid.UUID  `json:"id"`
	EntityType  string     `json:"entity_type"`
	EntityID    string     `json:"entity_id"`
	Kind        string     `json:"kind"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

// VisibleCondition hides drafts (no published_at) and future-dated articles everywhere in the API.
const VisibleCondition = "published_at IS NOT NULL AND published_at <= now()"

var Resource = crud.Resource[Article]{
	Path:    "/articles",
	Table:   "articles",
	Name:    "article",
	ID:      crud.IDUUID,
	Columns: []string{"id", "entity_type", "entity_id", "kind", "published_at", "created_at"},
	Scan: func(row database.Scanner) (Article, error) {
		var x Article
		err := row.Scan(&x.ID, &x.EntityType, &x.EntityID, &x.Kind, &x.PublishedAt, &x.CreatedAt)
		return x, err
	},
	Search: queries.EntitySearchConfig{
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
	},
}

// Swagger docs for the routes crud.Handler registers for Resource.

// @Summary Search articles
// @Description Retrieve a page of articles filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags articles
// @ID search-articles
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. published_at,asc)"
// @Success 200 {array} Article "List of matching articles"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /articles [get]
func _() {}

// @Summary Get an article
// @Tags articles
// @ID get-articles
// @Produce json
// @Param id path string true "Article UUID"
// @Success 200 {object} Article
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /articles/{id} [get]
func _() {}
