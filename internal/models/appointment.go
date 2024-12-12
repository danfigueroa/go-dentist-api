package models

import "fmt"

type Appointment struct {
	ID        string `json:"id"`
	DentistID string `json:"dentist_id"`
	PatientID string `json:"patient_id"`
	DateTime  string `json:"date_time"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// IsValid verifica se os campos obrigatórios do agendamento estão preenchidos
func (a *Appointment) IsValid() error {
	if a.DentistID == "" {
		return fmt.Errorf("dentist ID is required")
	}
	if a.PatientID == "" {
		return fmt.Errorf("patient ID is required")
	}
	if a.DateTime == "" {
		return fmt.Errorf("date and time is required")
	}

	return nil
}
