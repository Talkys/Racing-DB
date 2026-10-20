package driver

import (
	"github.com/google/uuid"
	"time"
)

type Driver struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Code        string    `json:"code"`
	Nationality string    `json:"nationality"`
	CreatedAt   time.Time `json:"created_at"`
}
