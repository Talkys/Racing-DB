package articlesearch

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"backend/modules/betterlogger"
	"backend/utils/apperror"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

const maxQueryRunes = 200

// Same rule as article_texts.language: a lowercase BCP 47 tag.
var languagePattern = regexp.MustCompile(`^[a-z]{2,3}(-[a-z0-9]{2,8})*$`)

type ArticleSearchHandler struct {
	DB      *sql.DB
	service *ArticleSearchService
	log     *betterlogger.Logger
}

func NewArticleSearchHandler(db *sql.DB) *ArticleSearchHandler {
	return &ArticleSearchHandler{DB: db, service: NewArticleSearchService(db), log: betterlogger.GetInstance()}
}

// RegisterRoutes must not clash with GET /articles/{id}: a literal segment is more specific, so it wins.
func (h *ArticleSearchHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/articles/search",
		),
		h.SearchArticles)
}

// SearchArticles finds articles by the words in their texts.
// @Summary Search articles
// @Description Full-text search over article titles, summaries and bodies: whole words, case-insensitive, no stemming. Supports "quoted phrases", or, and -exclusion. One result per article, using the text in `lang` when it qualifies. Without `q`, lists every article newest first. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags articles
// @ID search-articles-text
// @Produce json
// @Param q query string false "Search words (max 200 characters); omit to browse"
// @Param lang query string false "Preferred text language, e.g. pt-br"
// @Param page query int false "Page number for pagination (default: 1)"
// @Success 200 {array} ArticleSearchResult
// @Failure 400 {object} routes.ErrorResponse "Invalid q, lang or page"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /articles/search [get]
func (h *ArticleSearchHandler) SearchArticles(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	params := r.URL.Query()

	page, err := queries.ParsePage(params)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	q := strings.TrimSpace(params.Get("q"))
	if utf8.RuneCountInString(q) > maxQueryRunes {
		routes.WriteError(w, apperror.BadRequest("invalid q: longer than %d characters", maxQueryRunes))
		return
	}

	lang := strings.ToLower(strings.TrimSpace(params.Get("lang")))
	if lang != "" && !languagePattern.MatchString(lang) {
		routes.WriteError(w, apperror.BadRequest("invalid lang: '%s'", params.Get("lang")))
		return
	}

	results, err := h.service.SearchArticles(r.Context(), SearchCriteria{Query: q, Language: lang, Page: page})
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}
