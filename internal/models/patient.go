package models

import "fmt"

type Patient struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	DateOfBirth  string `json:"date_of_birth"`
	MedicalNotes string `json:"medical_notes"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// IsValid verifica se os campos obrigatórios do paciente estão preenchidos
func (p *Patient) IsValid() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Email == "" {
		return fmt.Errorf("email is required")
	}

	return nil
}
