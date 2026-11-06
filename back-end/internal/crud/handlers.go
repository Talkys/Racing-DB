package crud

import (
	"database/sql"
	"net/http"

	"backend/utils/constants"
	"backend/utils/queries"
	"backend/utils/routes"
)

// Handler serves GET Path and GET Path/{id}. Swagger annotations live next to each Resource,
// since swag cannot read them from a generic handler.
type Handler[T any] struct {
	res     *Resource[T]
	service *Service[T]
}

func NewHandler[T any](db *sql.DB, res Resource[T]) *Handler[T] {
	return &Handler[T]{res: &res, service: NewService(db, &res)}
}

func (h *Handler[T]) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(routes.BuildRoute("GET", constants.BaseApiUrl_V1, h.res.Path), h.Search)
	mux.HandleFunc(routes.BuildRoute("GET", constants.BaseApiUrl_V1, h.res.Path+"/{id}"), h.Get)
}

func (h *Handler[T]) Search(w http.ResponseWriter, r *http.Request) {
	sc, err := queries.BuildCriteriaFromURL(r.URL.Query())
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	results, err := h.service.Search(r.Context(), *sc)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, results)
}

func (h *Handler[T]) Get(w http.ResponseWriter, r *http.Request) {
	id, err := h.pathID(r)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	result, err := h.service.Get(r.Context(), id)
	if err != nil {
		routes.WriteError(w, err)
		return
	}

	routes.WriteJSONResponse(w, http.StatusOK, result)
}

func (h *Handler[T]) pathID(r *http.Request) (any, error) {
	if h.res.ID == IDSlug {
		return routes.SlugPathValue(r, "id")
	}
	return routes.UUIDPathValue(r, "id")
}
