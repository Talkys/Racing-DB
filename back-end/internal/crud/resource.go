// Package crud serves the read-only list and get-by-id endpoints shared by every plain entity.
// An entity describes itself once as a Resource; NewHandler builds the same handler → service → core →
// repository chain each entity used to write out by hand.
package crud

import (
	"backend/modules/database"
	"backend/utils/queries"
	"context"
)

// IDKind is the type of a resource's id path segment.
type IDKind int

const (
	IDUUID IDKind = iota // /races/{uuid}
	IDSlug               // /categories/{slug}, matched case-insensitively
)

// Resource is the whole configuration of one entity endpoint.
type Resource[T any] struct {
	Path    string // route under the API base, e.g. "/races"; the get route is Path + "/{id}"
	Table   string
	Name    string // singular, for messages: "race '…' not found"
	ID      IDKind
	Columns []string                          // selected in this order
	Scan    func(database.Scanner) (T, error) // reads Columns, in the same order
	Search  queries.EntitySearchConfig

	// Check runs in the core on every row before it is returned, for validation the database cannot do
	// (e.g. JSONB shapes). It may fix up the row in place; an error fails the request. Nil skips it.
	Check func(ctx context.Context, x *T) error
}
