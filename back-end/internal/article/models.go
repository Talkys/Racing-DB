package article

import (
	"github.com/google/uuid"
	"time"
)

type Article struct {
	ID          uuid.UUID  `json:"id"`
	EntityType  string     `json:"entity_type"`
	EntityID    string     `json:"entity_id"`
	Kind        string     `json:"kind"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
