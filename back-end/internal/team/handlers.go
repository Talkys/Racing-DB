package team

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type TeamHandler struct {
	DB      *sql.DB
	service *TeamService
	log     *betterlogger.Logger
}

func NewTeamHandler(db *sql.DB) *TeamHandler {
	return &TeamHandler{DB: db, service: NewTeamService(db), log: betterlogger.GetInstance()}
}

func (h *TeamHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/teams",
		),
		h.SearchTeams)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/teams/{id}",
		),
		h.GetTeam)
}

// SearchTeams searches teams based on URL query criteria.
// @Summary Search teams
// @Description Retrieve a page of teams filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags teams
// @ID search-teams
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Team "List of matching teams"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /teams [get]
func (h *TeamHandler) SearchTeams(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchTeams(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetTeam returns one team by id.
// @Summary Get a team
// @Tags teams
// @ID get-teams
// @Produce json
// @Param id path string true "Team UUID"
// @Success 200 {object} Team
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /teams/{id} [get]
func (h *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetTeam(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
