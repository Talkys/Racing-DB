package circuit

import (
	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
	"database/sql"
	"net/http"
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
			"/circuit",
		),
		h.SearchCircuits)
}

func (h *CircuitHandler) SearchCircuits(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	circuit, err := h.service.SearchCircuits(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, circuit)
}
