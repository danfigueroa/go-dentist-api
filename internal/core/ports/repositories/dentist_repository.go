package repositories

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type DentistRepository interface {
	Create(ctx context.Context, dentist *domain.Dentist) error
	GetByID(ctx context.Context, id string) (*domain.Dentist, error)
	GetByCRO(ctx context.Context, cro string) (*domain.Dentist, error)
	GetByName(ctx context.Context, name string) ([]*domain.Dentist, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Dentist, error)
	Update(ctx context.Context, dentist *domain.Dentist) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}) ([]*domain.Dentist, error)
}
