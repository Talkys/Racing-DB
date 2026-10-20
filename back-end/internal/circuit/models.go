package circuit

import (
	"time"

	"github.com/google/uuid"
)

type Circuit struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Country   string    `json:"country"`
	LengthKm  float32   `json:"length_km"`
	CreatedAt time.Time `json:"created_at"`
}
