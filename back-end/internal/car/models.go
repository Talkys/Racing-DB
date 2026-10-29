package car

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type Car struct {
	ID         uuid.UUID       `json:"id"`
	CategoryID string          `json:"category_id"`
	TeamID     uuid.UUID       `json:"team_id"`
	EngineID   uuid.UUID       `json:"engine_id"`
	ChassisID  uuid.UUID       `json:"chassis_id"`
	Name       string          `json:"name"`
	Season     int             `json:"season"`
	Parts      json.RawMessage `json:"parts" swaggertype:"object"`
	CreatedAt  time.Time       `json:"created_at"`
}
