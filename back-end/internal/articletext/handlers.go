package articletext

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type ArticleTextHandler struct {
	DB      *sql.DB
	service *ArticleTextService
	log     *betterlogger.Logger
}

func NewArticleTextHandler(db *sql.DB) *ArticleTextHandler {
	return &ArticleTextHandler{DB: db, service: NewArticleTextService(db), log: betterlogger.GetInstance()}
}

func (h *ArticleTextHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/article-texts",
		),
		h.SearchArticleTexts)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/article-texts/{id}",
		),
		h.GetArticleText)
}

// SearchArticleTexts searches article texts based on URL query criteria.
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
func (h *ArticleTextHandler) SearchArticleTexts(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchArticleTexts(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetArticleText returns one article text by id.
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
func (h *ArticleTextHandler) GetArticleText(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetArticleText(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
