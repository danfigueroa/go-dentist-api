package domain

import (
	"errors"
	"time"
)

type ProcedureStatus string

const (
	ProcedureStatusPending   ProcedureStatus = "pending"
	ProcedureStatusCompleted ProcedureStatus = "completed"
	ProcedureStatusCancelled ProcedureStatus = "cancelled"
)

type Procedure struct {
	ID            string          `json:"id"`
	ClinicID      string          `json:"clinic_id"`
	AppointmentID string          `json:"appointment_id"`
	Name          string          `json:"name"`
	Description   string          `json:"description"`
	Price         float64         `json:"price"`
	Status        ProcedureStatus `json:"status"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func (p *Procedure) Validate() error {
	if p.Name == "" {
		return errors.New("procedure name is required")
	}
	if p.ClinicID == "" {
		return errors.New("clinic ID is required")
	}
	if p.Price < 0 {
		return errors.New("price cannot be negative")
	}
	return nil
}

func NewProcedure(clinicID, appointmentID, name, description string, price float64) *Procedure {
	now := time.Now()
	return &Procedure{
		ClinicID:      clinicID,
		AppointmentID: appointmentID,
		Name:          name,
		Description:   description,
		Price:         price,
		Status:        ProcedureStatusPending,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}
