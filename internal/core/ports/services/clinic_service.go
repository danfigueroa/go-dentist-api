package services

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

type CreateClinicRequest struct {
	Name    string `json:"name"`
	CNPJ    string `json:"cnpj"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type UpdateClinicRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
}

type ClinicService interface {
	CreateClinic(ctx context.Context, req CreateClinicRequest) (*domain.Clinic, error)
	GetClinic(ctx context.Context, id string) (*domain.Clinic, error)
	GetClinicByCNPJ(ctx context.Context, cnpj string) (*domain.Clinic, error)
	UpdateClinic(ctx context.Context, id string, req UpdateClinicRequest) error
	DeleteClinic(ctx context.Context, id string) error
	ListClinics(ctx context.Context, filters map[string]interface{}) ([]*domain.Clinic, error)
}
