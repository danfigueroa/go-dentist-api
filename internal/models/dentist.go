package models

import "fmt"

type Dentist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CRO       string `json:"cro"`
	Country   string `json:"country"`
	Specialty string `json:"specialty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (d *Dentist) IsValid() error {
	fields := map[string]string{
		"id":      d.ID,
		"name":    d.Name,
		"email":   d.Email,
		"CRO":     d.CRO,
		"country": d.Country,
	}

	for field, value := range fields {
		if value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}

	return nil
}
