package services

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type CreateProcedureRequest struct {
	ClinicID      string  `json:"clinic_id"`
	AppointmentID string  `json:"appointment_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
}

type UpdateProcedureRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ProcedureService interface {
	CreateProcedure(ctx context.Context, req CreateProcedureRequest) (*domain.Procedure, error)
	GetProcedure(ctx context.Context, id string) (*domain.Procedure, error)
	GetProceduresByAppointment(ctx context.Context, appointmentID string) ([]*domain.Procedure, error)
	GetProceduresByClinic(ctx context.Context, clinicID string) ([]*domain.Procedure, error)
	GetProceduresByName(ctx context.Context, name string) ([]*domain.Procedure, error)
	UpdateProcedure(ctx context.Context, id string, req UpdateProcedureRequest) error
	DeleteProcedure(ctx context.Context, id string) error
	ListProcedures(ctx context.Context, filters map[string]interface{}) ([]*domain.Procedure, error)
}
