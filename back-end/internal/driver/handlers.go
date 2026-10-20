package driver

import (
	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
	"database/sql"
	"net/http"
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
			"/driver",
		),
		h.SearchDrivers)
}

func (h *DriverHandler) SearchDrivers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := h.service.SearchDrivers(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}
