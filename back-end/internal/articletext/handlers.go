package articletext

import (
	"backend/modules/betterlogger"
	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
	"database/sql"
	"net/http"
)

type ArticleTextHandler struct {
	DB      *sql.DB
	service *ArticleTextService
	log     *betterlogger.Logger
}

func NewArticleTextHandler(db *sql.DB) *ArticleTextHandler {
	return &ArticleTextHandler{DB: db, service: NewArticleTextService(db), log: betterlogger.GetInstance()}
}

func (h *ArticleTextHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		routes.BuildRoute(
			"GET",
			constants.BaseApiUrl_V1,
			"/article-text",
		),
		h.SearchArticleTexts)
}

func (h *ArticleTextHandler) SearchArticleTexts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	query := r.URL.Query()
	h.log.Info("%s", r.URL.String())
	sc, err := queries.BuildCriteriaFromURL(query)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := h.service.SearchArticleTexts(ctx, *sc)
	if err != nil {
		h.log.Error("%s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}
