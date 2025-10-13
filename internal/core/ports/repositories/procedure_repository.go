package repositories

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type ProcedureRepository interface {
	Create(ctx context.Context, procedure *domain.Procedure) error
	GetByID(ctx context.Context, id string) (*domain.Procedure, error)
	GetByAppointmentID(ctx context.Context, appointmentID string) ([]*domain.Procedure, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Procedure, error)
	GetByName(ctx context.Context, name string) ([]*domain.Procedure, error)
	Update(ctx context.Context, procedure *domain.Procedure) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}) ([]*domain.Procedure, error)
}
