package usecases

import (
	"context"
	"errors"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports/repositories"
	"go-dentist-api/internal/core/ports/services"
)

type ClinicUseCase struct {
	clinicRepo repositories.ClinicRepository
}

func NewClinicUseCase(clinicRepo repositories.ClinicRepository) services.ClinicService {
	return &ClinicUseCase{
		clinicRepo: clinicRepo,
	}
}

func (uc *ClinicUseCase) CreateClinic(ctx context.Context, req services.CreateClinicRequest) (*domain.Clinic, error) {
	existingClinic, _ := uc.clinicRepo.GetByCNPJ(ctx, req.CNPJ)
	if existingClinic != nil {
		return nil, errors.New("clinic with this CNPJ already exists")
	}

	clinic, err := domain.NewClinic(req.Name, req.CNPJ, req.Address, req.Phone, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid clinic data: %w", err)
	}

	if err := uc.clinicRepo.Create(ctx, clinic); err != nil {
		return nil, fmt.Errorf("failed to create clinic: %w", err)
	}

	return clinic, nil
}

func (uc *ClinicUseCase) GetClinic(ctx context.Context, id string) (*domain.Clinic, error) {
	clinic, err := uc.clinicRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinic: %w", err)
	}

	return clinic, nil
}

func (uc *ClinicUseCase) GetClinicByCNPJ(ctx context.Context, cnpj string) (*domain.Clinic, error) {
	clinic, err := uc.clinicRepo.GetByCNPJ(ctx, cnpj)
	if err != nil {
		return nil, fmt.Errorf("failed to get clinic by CNPJ: %w", err)
	}

	return clinic, nil
}

func (uc *ClinicUseCase) UpdateClinic(ctx context.Context, id string, req services.UpdateClinicRequest) error {
	clinic, err := uc.clinicRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}

	if req.Name != "" {
		clinic.Name = req.Name
	}
	if req.Address != "" {
		clinic.Address = req.Address
	}
	if req.Phone != "" {
		clinic.Phone = req.Phone
	}
	if req.Email != "" {
		clinic.Email = req.Email
	}

	if err := clinic.Validate(); err != nil {
		return fmt.Errorf("invalid clinic data: %w", err)
	}

	if err := uc.clinicRepo.Update(ctx, clinic); err != nil {
		return fmt.Errorf("failed to update clinic: %w", err)
	}

	return nil
}

func (uc *ClinicUseCase) DeleteClinic(ctx context.Context, id string) error {
	_, err := uc.clinicRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("clinic not found: %w", err)
	}

	if err := uc.clinicRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete clinic: %w", err)
	}

	return nil
}

func (uc *ClinicUseCase) ListClinics(ctx context.Context, filters map[string]interface{}) ([]*domain.Clinic, error) {
	clinics, err := uc.clinicRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list clinics: %w", err)
	}

	return clinics, nil
}
