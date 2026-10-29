package article

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type ArticleHandler struct {
	DB      *sql.DB
	service *ArticleService
	log     *betterlogger.Logger
}

func NewArticleHandler(db *sql.DB) *ArticleHandler {
	return &ArticleHandler{DB: db, service: NewArticleService(db), log: betterlogger.GetInstance()}
}

func (h *ArticleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/articles",
		),
		h.SearchArticles)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/articles/{id}",
		),
		h.GetArticle)
}

// SearchArticles searches articles based on URL query criteria.
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
func (h *ArticleHandler) SearchArticles(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchArticles(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetArticle returns one article by id.
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
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetArticle(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
