package articlesearch

import (
	"time"

	"github.com/google/uuid"
)

// PreviewSegment is a piece of a result's preview. Match marks the words that matched the query,
// so clients can highlight them without rendering HTML.
type PreviewSegment struct {
	Text  string `json:"text"`
	Match bool   `json:"match"`
}

// ArticleSearchResult is one article with the text chosen to represent it: the requested language
// when that text qualifies, otherwise the best match. Language, Title and Summary are null only when
// browsing (no query) an article that has no texts.
type ArticleSearchResult struct {
	ID          uuid.UUID        `json:"id"`
	EntityType  string           `json:"entity_type"`
	EntityID    string           `json:"entity_id"`
	Kind        string           `json:"kind"`
	PublishedAt *time.Time       `json:"published_at"`
	Language    *string          `json:"language"`
	Title       *string          `json:"title"`
	Summary     *string          `json:"summary"`
	Preview     []PreviewSegment `json:"preview"`
}

type SearchCriteria struct {
	Query    string // empty = browse every article, newest first
	Language string // preferred text language; may be empty
	Page     int
}
