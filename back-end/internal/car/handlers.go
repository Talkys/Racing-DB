package car

import (
	"database/sql"
	"net/http"

	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

type CarHandler struct {
	DB      *sql.DB
	service *CarService
	log     *betterlogger.Logger
}

func NewCarHandler(db *sql.DB) *CarHandler {
	return &CarHandler{DB: db, service: NewCarService(db), log: betterlogger.GetInstance()}
}

func (h *CarHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/cars",
		),
		h.SearchCars)
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/cars/{id}",
		),
		h.GetCar)
}

// SearchCars searches cars based on URL query criteria.
// @Summary Search cars
// @Description Retrieve a page of cars filtered and sorted by query criteria. Returns up to 26 rows: the 26th only signals that a next page exists.
// @Tags cars
// @ID search-cars
// @Produce json
// @Param page query int false "Page number for pagination (default: 1)"
// @Param filter query string false "Filter expression in format: column,operator,value (operators: eq neq gt gte lt lte ilike in; in takes values separated by |)"
// @Param sort query string false "Sort expression in format: column,direction (e.g. name,asc)"
// @Success 200 {array} Car "List of matching cars"
// @Failure 400 {object} routes.ErrorResponse "Invalid page, filter or sort"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /cars [get]
func (h *CarHandler) SearchCars(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.SearchCars(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

// GetCar returns one car by id.
// @Summary Get a car
// @Tags cars
// @ID get-cars
// @Produce json
// @Param id path string true "Car UUID"
// @Success 200 {object} Car
// @Failure 400 {object} routes.ErrorResponse "Invalid id"
// @Failure 404 {object} routes.ErrorResponse "Not found"
// @Failure 500 {object} routes.ErrorResponse "Internal Server Error"
// @Router /cars/{id} [get]
func (h *CarHandler) GetCar(w http.ResponseWriter, r *http.Request) {
	h.log.Info("%s", r.URL.String())
	id, err := routes.UUIDPathValue(r, "id")
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.GetCar(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}
