package driver

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type DriverHandler struct {
	DB      *sql.DB
	service *DriverService
	log     *betterlogger.Logger
}

func NewDriverHandler(db *sql.DB) *DriverHandler {
	return &DriverHandler{DB: db, service: NewDriverService(db), log: betterlogger.GetInstance()}
}

func (h *DriverHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/drivers",
		),
		h.SearchDrivers)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/drivers/{id}",
		),
		h.GetDriver)
}

// SearchDrivers searches drivers based on URL query criteria.
// @Summary Search drivers
// @Description Retrieve a page of drivers filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags drivers
// @ID search-drivers
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. last_name,asc)"
// @Success 200 {array} Driver "List of matching drivers"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /drivers [get]
func (h *DriverHandler) SearchDrivers(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchDrivers(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetDriver returns one driver by id.
// @Summary Get a driver
// @Tags drivers
// @ID get-drivers
// @Produce json
// @Param id path string true "Driver UUID"
// @Success 200 {object} Driver
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /drivers/{id} [get]
func (h *DriverHandler) GetDriver(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetDriver(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
