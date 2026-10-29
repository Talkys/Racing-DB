package circuit

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type CircuitHandler struct {
	DB      *sql.DB
	service *CircuitService
	log     *betterlogger.Logger
}

func NewCircuitHandler(db *sql.DB) *CircuitHandler {
	return &CircuitHandler{DB: db, service: NewCircuitService(db), log: betterlogger.GetInstance()}
}

func (h *CircuitHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/circuits",
		),
		h.SearchCircuits)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/circuits/{id}",
		),
		h.GetCircuit)
}

// SearchCircuits searches circuits based on URL query criteria.
// @Summary Search circuits
// @Description Retrieve a page of circuits filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags circuits
// @ID search-circuits
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Circuit "List of matching circuits"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /circuits [get]
func (h *CircuitHandler) SearchCircuits(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchCircuits(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetCircuit returns one circuit by id.
// @Summary Get a circuit
// @Tags circuits
// @ID get-circuits
// @Produce json
// @Param id path string true "Circuit UUID"
// @Success 200 {object} Circuit
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /circuits/{id} [get]
func (h *CircuitHandler) GetCircuit(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetCircuit(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
