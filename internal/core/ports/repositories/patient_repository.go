package repositories

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type PatientRepository interface {
	Create(ctx context.Context, patient *domain.Patient) error
	GetByID(ctx context.Context, id string) (*domain.Patient, error)
	GetByCPF(ctx context.Context, cpf string) (*domain.Patient, error)
	GetByName(ctx context.Context, name string) ([]*domain.Patient, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error)
	Update(ctx context.Context, patient *domain.Patient) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}) ([]*domain.Patient, error)
}
