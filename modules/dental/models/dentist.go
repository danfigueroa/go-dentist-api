package models

import (
	"fmt"
	"time"
)

type Dentist struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CRO       string    `json:"cro"`
	Country   string    `json:"country"`
	Specialty string    `json:"specialty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Dentist) IsValid() error {
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if d.Email == "" {
		return fmt.Errorf("email is required")
	}
	if d.CRO == "" {
		return fmt.Errorf("CRO is required")
	}
	if d.Country == "" {
		return fmt.Errorf("country is required")
	}

	return nil
}