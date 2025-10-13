package services

import (
	"context"
	"go-dentist-api/internal/core/domain"
	"time"
)

type CreateAppointmentRequest struct {
	ClinicID    string    `json:"clinic_id"`
	PatientID   string    `json:"patient_id"`
	DentistID   string    `json:"dentist_id"`
	UserID      string    `json:"user_id"`
	ScheduledAt time.Time `json:"scheduled_at"`
	Notes       string    `json:"notes"`
}

type UpdateAppointmentRequest struct {
	ScheduledAt time.Time `json:"scheduled_at"`
	Notes       string    `json:"notes"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount"`
}

type AppointmentService interface {
	CreateAppointment(ctx context.Context, req CreateAppointmentRequest) (*domain.Appointment, error)
	GetAppointment(ctx context.Context, id string) (*domain.Appointment, error)
	GetAppointmentsByPatient(ctx context.Context, patientID string) ([]*domain.Appointment, error)
	GetAppointmentsByDentist(ctx context.Context, dentistID string) ([]*domain.Appointment, error)
	GetAppointmentsByClinic(ctx context.Context, clinicID string) ([]*domain.Appointment, error)
	GetAppointmentsByDateRange(ctx context.Context, clinicID string, startDate, endDate time.Time) ([]*domain.Appointment, error)
	UpdateAppointment(ctx context.Context, id string, req UpdateAppointmentRequest) error
	DeleteAppointment(ctx context.Context, id string) error
	ListAppointments(ctx context.Context, filters map[string]interface{}) ([]*domain.Appointment, error)
}
