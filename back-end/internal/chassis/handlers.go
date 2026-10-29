package chassis

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type ChassisHandler struct {
	DB      *sql.DB
	service *ChassisService
	log     *betterlogger.Logger
}

func NewChassisHandler(db *sql.DB) *ChassisHandler {
	return &ChassisHandler{DB: db, service: NewChassisService(db), log: betterlogger.GetInstance()}
}

func (h *ChassisHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/chassis",
		),
		h.SearchChassis)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/chassis/{id}",
		),
		h.GetChassis)
}

// SearchChassis searches chassis based on URL query criteria.
// @Summary Search chassis
// @Description Retrieve a page of chassis filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags chassis
// @ID search-chassis
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Chassis "List of matching chassis"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /chassis [get]
func (h *ChassisHandler) SearchChassis(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchChassis(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetChassis returns one chassis by id.
// @Summary Get a chassis
// @Tags chassis
// @ID get-chassis
// @Produce json
// @Param id path string true "Chassis UUID"
// @Success 200 {object} Chassis
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /chassis/{id} [get]
func (h *ChassisHandler) GetChassis(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetChassis(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
