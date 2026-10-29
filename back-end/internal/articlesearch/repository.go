package articlesearch

import (
	"backend/internal/article"
	"backend/modules/database"
	"context"
	"fmt"
)

const PageSize = 25

// Highlight markers ts_headline puts around matched words (Unicode private use characters, which
// markdown_to_text strips from the source text so they cannot appear otherwise).
const (
	matchStart = ""
	matchStop  = ""
)

var headlineOptions = fmt.Sprintf(`StartSel=%s, StopSel=%s, MaxWords=35, MinWords=15, MaxFragments=2, FragmentDelimiter=" … "`, matchStart, matchStop)

// Per article, the text in the requested language wins, then one sharing its primary subtag (pt for pt-br).
const languagePreference = `CASE WHEN t.language = $%[1]d THEN 0
	WHEN split_part(t.language, '-', 1) = split_part($%[1]d, '-', 1) THEN 1 ELSE 2 END`

// searchQuery ranks matching texts, keeps one per article, paginates, and only then reads bodies and
// builds previews: ts_headline re-reads the whole text, so it runs on at most one page of rows. It reads
// body_text, the Markdown-free copy of body_md that Postgres maintains on write.
// Params: $1 query, $2 language, $3 limit, $4 offset, $5 headline options.
var searchQuery = `
WITH q AS (SELECT websearch_to_tsquery('simple', $1) AS query),
best AS (
	SELECT DISTINCT ON (t.article_id) t.article_id, t.id AS text_id, ts_rank_cd(t.search_vector, q.query) AS rank
	FROM article_texts t CROSS JOIN q
	WHERE t.search_vector @@ q.query
		AND t.article_id IN (SELECT id FROM articles WHERE ` + article.VisibleCondition + `)
	ORDER BY t.article_id, ` + fmt.Sprintf(languagePreference, 2) + `, rank DESC, t.language
),
page AS (
	SELECT a.id, a.entity_type, a.entity_id, a.kind, a.published_at, b.text_id, b.rank
	FROM best b JOIN articles a ON a.id = b.article_id
	ORDER BY b.rank DESC, a.published_at DESC, a.id
	LIMIT $3 OFFSET $4
)
SELECT p.id, p.entity_type, p.entity_id, p.kind, p.published_at, t.language, t.title, t.summary,
	ts_headline('simple', t.body_text, q.query, $5)
FROM page p JOIN article_texts t ON t.id = p.text_id CROSS JOIN q
ORDER BY p.rank DESC, p.published_at DESC, p.id`

// browseQuery lists every visible article newest first, with the start of its preferred text as preview.
// Params: $1 language, $2 limit, $3 offset.
var browseQuery = `
WITH best AS (
	SELECT DISTINCT ON (a.id) a.id, a.entity_type, a.entity_id, a.kind, a.published_at, t.id AS text_id
	FROM articles a LEFT JOIN article_texts t ON t.article_id = a.id
	WHERE ` + article.VisibleCondition + `
	ORDER BY a.id, ` + fmt.Sprintf(languagePreference, 1) + `, t.language
),
page AS (
	SELECT * FROM best
	ORDER BY published_at DESC, id
	LIMIT $2 OFFSET $3
)
SELECT p.id, p.entity_type, p.entity_id, p.kind, p.published_at, t.language, t.title, t.summary,
	left(t.body_text, ` + fmt.Sprint(browsePreviewSource) + `)
FROM page p LEFT JOIN article_texts t ON t.id = p.text_id
ORDER BY p.published_at DESC, p.id`

// browsePreviewSource is how much text a browse preview reads; plenty for previewRunes after trimming the title.
const browsePreviewSource = 1000

// rawResult is a row as the database returns it; Preview still holds the raw headline or text.
type rawResult struct {
	ArticleSearchResult
	Preview *string
}

type ArticleSearchRepository struct{}

func NewArticleSearchRepository() *ArticleSearchRepository {
	return &ArticleSearchRepository{}
}

// SearchArticles returns up to PageSize + 1 rows; the extra row signals a next page.
func (r *ArticleSearchRepository) SearchArticles(ctx context.Context, sc SearchCriteria) ([]rawResult, error) {
	exec := database.ExtractTx(ctx)
	if exec == nil {
		return nil, fmt.Errorf("No transaction provided")
	}

	limit := PageSize + 1
	offset := (sc.Page - 1) * PageSize

	query, args := browseQuery, []any{sc.Language, limit, offset}
	if sc.Query != "" {
		query, args = searchQuery, []any{sc.Query, sc.Language, limit, offset, headlineOptions}
	}

	rows, err := exec.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []rawResult{}

	for rows.Next() {
		var x rawResult
		if err := rows.Scan(&x.ID, &x.EntityType, &x.EntityID, &x.Kind, &x.PublishedAt, &x.Language, &x.Title, &x.Summary, &x.Preview); err != nil {
			return nil, err
		}
		results = append(results, x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
