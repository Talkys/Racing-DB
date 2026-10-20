package engine

import (
	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
	"database/sql"
	"net/http"
)

type EngineHandler struct {
	DB      *sql.DB
	service *EngineService
	log     *betterlogger.Logger
}

func NewEngineHandler(db *sql.DB) *EngineHandler {
	return &EngineHandler{DB: db, service: NewEngineService(db), log: betterlogger.GetInstance()}
}

func (h *EngineHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/engine",
		),
		h.SearchEngines)
}

func (h *EngineHandler) SearchEngines(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := h.service.SearchEngines(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}
