package race

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type RaceHandler struct {
	DB      *sql.DB
	service *RaceService
	log     *betterlogger.Logger
}

func NewRaceHandler(db *sql.DB) *RaceHandler {
	return &RaceHandler{DB: db, service: NewRaceService(db), log: betterlogger.GetInstance()}
}

func (h *RaceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/races",
		),
		h.SearchRaces)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/races/{id}",
		),
		h.GetRace)
}

// SearchRaces searches races based on URL query criteria.
// @Summary Search races
// @Description Retrieve a page of races filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags races
// @ID search-races
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. race_date,asc)"
// @Success 200 {array} Race "List of matching races"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /races [get]
func (h *RaceHandler) SearchRaces(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchRaces(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetRace returns one race by id.
// @Summary Get a race
// @Tags races
// @ID get-races
// @Produce json
// @Param id path string true "Race UUID"
// @Success 200 {object} Race
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /races/{id} [get]
func (h *RaceHandler) GetRace(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetRace(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
