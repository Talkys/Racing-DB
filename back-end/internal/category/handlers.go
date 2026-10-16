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
			"/active_categories",
		),
		h.ListCategories)

	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/categories",
		),
		h.SearchCategories)
}

func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	categories, err := h.service.ListCategories(ctx)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, categories)
}

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
