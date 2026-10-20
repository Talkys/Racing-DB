package engine

import (
	"github.com/google/uuid"
	"time"
)

type Engine struct {
	ID             uuid.UUID `json:"id"`
	Manufacturer   string    `json:"manufacturer"`
	Name           string    `json:"name"`
	Configuration  string    `json:"configuration"`
	DisplacementCC int       `json:"displacement_cc"`
	Cylinders      int16     `json:"cylinders"`
	Aspiration     string    `json:"aspiration"`
	PowerKW        *float64  `json:"power_kw"`
	CreatedAt      time.Time `json:"created_at"`
}
