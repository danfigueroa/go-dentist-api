package models

import "fmt"

type Procedure struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Duration    string `json:"duration"` // em minutos
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// IsValid verifica se os campos obrigatórios do procedimento estão preenchidos
func (p *Procedure) IsValid() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Price == "" {
		return fmt.Errorf("price is required")
	}
	if p.Duration == "" {
		return fmt.Errorf("duration is required")
	}

	return nil
}