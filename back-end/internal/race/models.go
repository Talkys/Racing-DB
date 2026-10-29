package race

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Race struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID string          `json:"category_id"`
	CircuitID  uuid.UUID       `json:"circuit_id"`
	WinnerID   *uuid.UUID      `json:"winner_id"`
	Name       string          `json:"name"`
	Season     int             `json:"season"`
	Round      int             `json:"round"`
	RaceDate   time.Time       `json:"race_date"`
	Payload    json.RawMessage `json:"payload" swaggertype:"object"`
	CreatedAt  time.Time       `json:"created_at"`
}
