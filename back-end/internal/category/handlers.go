package category

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type CategoryHandler struct {
	DB      *sql.DB
	service *CategoryService
	log     *betterlogger.Logger
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{DB: db, service: NewCategoryService(db), log: betterlogger.GetInstance()}
}

func (h *CategoryHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/categories",
		),
		h.SearchCategories)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/categories/{id}",
		),
		h.GetCategory)
}

// SearchCategories searches categories based on URL query criteria.
// @Summary Search categories
// @Description Retrieve a page of categories filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags categories
// @ID search-categories
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Category "List of matching categories"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /categories [get]
func (h *CategoryHandler) SearchCategories(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchCategories(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetCategory returns one category by id.
// @Summary Get a category
// @Tags categories
// @ID get-categories
// @Produce json
// @Param id path string true "Category slug, e.g. f1 (case-insensitive)"
// @Success 200 {object} Category
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.SlugPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetCategory(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
