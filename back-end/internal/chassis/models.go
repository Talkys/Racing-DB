package chassis

import (
	"github.com/google/uuid"
	"time"
)

type Chassis struct {
	ID           uuid.UUID `json:"id"`
	Manufacturer string    `json:"manufacturer"`
	Name         string    `json:"name"`
	Material     string    `json:"material"`
	WheelbaseMM  *int      `json:"wheelbase_mm"`
	WeightKG     *float64  `json:"weight_kg"`
	CreatedAt    time.Time `json:"created_at"`
}
