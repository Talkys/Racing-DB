package routes

import (
	"backend/utils/apperror"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9_-]+$`)

// UUIDPathValue reads a UUID path segment, e.g. {id} in /races/{id}.
func UUIDPathValue(r *http.Request, name string) (uuid.UUID, error) {
	raw := r.PathValue(name)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperror.BadRequest("invalid %s: '%s' is not a UUID", name, raw)
	}
	return id, nil
}

// SlugPathValue reads a lowercase slug path segment (category ids). Input is lowercased first, so F1 == f1.
func SlugPathValue(r *http.Request, name string) (string, error) {
	slug := strings.ToLower(strings.TrimSpace(r.PathValue(name)))
	if !slugPattern.MatchString(slug) {
		return "", apperror.BadRequest("invalid %s: '%s' is not a valid slug", name, r.PathValue(name))
	}
	return slug, nil
}
