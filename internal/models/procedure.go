package models

type Procedure struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Duration    string `json:"duration"` // em minutos
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
