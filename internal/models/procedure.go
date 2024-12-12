package models

import (
	"fmt"
	"time"
)

type Procedure struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	DentistID    string    `json:"dentist_id"`
	PatientID    string    `json:"patient_id"`
	PerformedAt  time.Time `json:"performed_at"`
	Observations string    `json:"observations"`
	Cost         float64   `json:"cost"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (p *Procedure) IsValid() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}
	if p.Type == "" {
		return fmt.Errorf("type is required")
	}
	if p.DentistID == "" {
		return fmt.Errorf("dentist_id is required")
	}
	if p.PerformedAt.IsZero() {
		return fmt.Errorf("performed_at is required")
	}
	return nil
}
