package car

import (
	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
	"database/sql"
	"net/http"
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
			"/car",
		),
		h.SearchCars)
}

func (h *CarHandler) SearchCars(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := h.service.SearchCars(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}
