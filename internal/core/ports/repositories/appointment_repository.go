package repositories

import (
	"context"
	"go-dentist-api/internal/core/domain"
	"time"
)

type AppointmentRepository interface {
	Create(ctx context.Context, appointment *domain.Appointment) error
	GetByID(ctx context.Context, id string) (*domain.Appointment, error)
	GetByPatientID(ctx context.Context, patientID string) ([]*domain.Appointment, error)
	GetByDentistID(ctx context.Context, dentistID string) ([]*domain.Appointment, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Appointment, error)
	GetByDateRange(ctx context.Context, clinicID string, startDate, endDate time.Time) ([]*domain.Appointment, error)
	Update(ctx context.Context, appointment *domain.Appointment) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}) ([]*domain.Appointment, error)
}
