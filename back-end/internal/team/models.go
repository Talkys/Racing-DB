package team

import (
	"github.com/google/uuid"
	"time"
)

type Team struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	ShortName   string    `json:"short_name"`
	Country     string    `json:"country"`
	FoundedYear *int      `json:"founded_year"`
	CreatedAt   time.Time `json:"created_at"`
}
