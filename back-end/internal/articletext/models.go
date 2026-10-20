package articletext

import (
	"github.com/google/uuid"
	"time"
)

type ArticleText struct {
	ID        uuid.UUID `json:"id"`
	ArticleID uuid.UUID `json:"article_id"`
	Language  string    `json:"language"`
	Title     string    `json:"title"`
	Summary   *string   `json:"summary"`
	BodyMd    string    `json:"body_md"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
