package repositories

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type ClinicRepository interface {
	Create(ctx context.Context, clinic *domain.Clinic) error
	GetByID(ctx context.Context, id string) (*domain.Clinic, error)
	GetByCNPJ(ctx context.Context, cnpj string) (*domain.Clinic, error)
	Update(ctx context.Context, clinic *domain.Clinic) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filters map[string]interface{}) ([]*domain.Clinic, error)
}
