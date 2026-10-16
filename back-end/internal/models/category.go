package models

type Category struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Is_active   bool   `json:"is_active"`
}
