package engine

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type EngineHandler struct {
	DB      *sql.DB
	service *EngineService
	log     *betterlogger.Logger
}

func NewEngineHandler(db *sql.DB) *EngineHandler {
	return &EngineHandler{DB: db, service: NewEngineService(db), log: betterlogger.GetInstance()}
}

func (h *EngineHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/engines",
		),
		h.SearchEngines)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/engines/{id}",
		),
		h.GetEngine)
}

// SearchEngines searches engines based on URL query criteria.
// @Summary Search engines
// @Description Retrieve a page of engines filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags engines
// @ID search-engines
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Engine "List of matching engines"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /engines [get]
func (h *EngineHandler) SearchEngines(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchEngines(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetEngine returns one engine by id.
// @Summary Get an engine
// @Tags engines
// @ID get-engines
// @Produce json
// @Param id path string true "Engine UUID"
// @Success 200 {object} Engine
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /engines/{id} [get]
func (h *EngineHandler) GetEngine(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetEngine(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
