package chassis

import (
	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
	"database/sql"
	"net/http"
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
		h.SearchChassiss)
}

func (h *ChassisHandler) SearchChassiss(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := h.service.SearchChassiss(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}
