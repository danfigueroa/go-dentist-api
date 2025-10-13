package services

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type CreateDentistRequest struct {
	ClinicID  string `json:"clinic_id"`
	Name      string `json:"name"`
	CRO       string `json:"cro"`
	Specialty string `json:"specialty"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type UpdateDentistRequest struct {
	Name      string `json:"name"`
	Specialty string `json:"specialty"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Active    *bool  `json:"active"`
}

type DentistService interface {
	CreateDentist(ctx context.Context, req CreateDentistRequest) (*domain.Dentist, error)
	GetDentist(ctx context.Context, id string) (*domain.Dentist, error)
	GetDentistByCRO(ctx context.Context, cro string) (*domain.Dentist, error)
	GetDentistsByName(ctx context.Context, name string) ([]*domain.Dentist, error)
	GetDentistsByClinic(ctx context.Context, clinicID string) ([]*domain.Dentist, error)
	UpdateDentist(ctx context.Context, id string, req UpdateDentistRequest) error
	DeleteDentist(ctx context.Context, id string) error
	ListDentists(ctx context.Context, filters map[string]interface{}) ([]*domain.Dentist, error)
}
