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
}

// SearchCategories searches categories based on URL query criteria.
// @Summary Search categories
// @Description Retrieve a list of categories filtered, sorted, and paginated by query criteria
// @ID search-categories
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (e.g. name,eq,electronics)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Category "List of matching categories"
// @Failure 500 {string} string "Internal Server Error"
// @Router /categories [get]
func (h *CategoryHandler) SearchCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := h.service.SearchCategories(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, categories)
}
