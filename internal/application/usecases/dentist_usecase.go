package usecases

import (
	"context"
	"errors"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports/repositories"
	"go-dentist-api/internal/core/ports/services"
)

type DentistUseCase struct {
	dentistRepo repositories.DentistRepository
	clinicRepo  repositories.ClinicRepository
}

func NewDentistUseCase(dentistRepo repositories.DentistRepository, clinicRepo repositories.ClinicRepository) services.DentistService {
	return &DentistUseCase{
		dentistRepo: dentistRepo,
		clinicRepo:  clinicRepo,
	}
}

func (uc *DentistUseCase) CreateDentist(ctx context.Context, req services.CreateDentistRequest) (*domain.Dentist, error) {
	clinic, err := uc.clinicRepo.GetByID(ctx, req.ClinicID)
	if err != nil {
		return nil, fmt.Errorf("clinic not found: %w", err)
	}
	if clinic == nil {
		return nil, errors.New("clinic not found")
	}

	existingDentist, _ := uc.dentistRepo.GetByCRO(ctx, req.CRO)
	if existingDentist != nil {
		return nil, errors.New("dentist with this CRO already exists")
	}

	dentist, err := domain.NewDentist(req.ClinicID, req.Name, req.CRO, req.Specialty, req.Phone, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid dentist data: %w", err)
	}

	if err := uc.dentistRepo.Create(ctx, dentist); err != nil {
		return nil, fmt.Errorf("failed to create dentist: %w", err)
	}

	return dentist, nil
}

func (uc *DentistUseCase) GetDentist(ctx context.Context, id string) (*domain.Dentist, error) {
	if id == "" {
		return nil, errors.New("dentist ID is required")
	}

	dentist, err := uc.dentistRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dentist: %w", err)
	}

	return dentist, nil
}

func (uc *DentistUseCase) GetDentistByCRO(ctx context.Context, cro string) (*domain.Dentist, error) {
	if cro == "" {
		return nil, errors.New("CRO is required")
	}

	dentist, err := uc.dentistRepo.GetByCRO(ctx, cro)
	if err != nil {
		return nil, fmt.Errorf("failed to get dentist by CRO: %w", err)
	}

	return dentist, nil
}

func (uc *DentistUseCase) GetDentistsByName(ctx context.Context, name string) ([]*domain.Dentist, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	dentists, err := uc.dentistRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get dentists by name: %w", err)
	}

	return dentists, nil
}

func (uc *DentistUseCase) GetDentistsByClinic(ctx context.Context, clinicID string) ([]*domain.Dentist, error) {
	if clinicID == "" {
		return nil, errors.New("clinic ID is required")
	}

	dentists, err := uc.dentistRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dentists by clinic: %w", err)
	}

	return dentists, nil
}

func (uc *DentistUseCase) UpdateDentist(ctx context.Context, id string, req services.UpdateDentistRequest) error {
	if id == "" {
		return errors.New("dentist ID is required")
	}

	dentist, err := uc.dentistRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get dentist: %w", err)
	}
	if dentist == nil {
		return errors.New("dentist not found")
	}

	if req.Name != "" {
		dentist.Name = req.Name
	}
	if req.Specialty != "" {
		dentist.Specialty = req.Specialty
	}
	if req.Phone != "" {
		dentist.Phone = req.Phone
	}
	if req.Email != "" {
		dentist.Email = req.Email
	}
	if req.Active != nil {
		dentist.Active = *req.Active
	}

	if err := dentist.Validate(); err != nil {
		return fmt.Errorf("invalid dentist data: %w", err)
	}

	if err := uc.dentistRepo.Update(ctx, dentist); err != nil {
		return fmt.Errorf("failed to update dentist: %w", err)
	}

	return nil
}

func (uc *DentistUseCase) DeleteDentist(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("dentist ID is required")
	}

	dentist, err := uc.dentistRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get dentist: %w", err)
	}
	if dentist == nil {
		return errors.New("dentist not found")
	}

	if err := uc.dentistRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete dentist: %w", err)
	}

	return nil
}

func (uc *DentistUseCase) ListDentists(ctx context.Context, filters map[string]interface{}) ([]*domain.Dentist, error) {
	dentists, err := uc.dentistRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list dentists: %w", err)
	}

	return dentists, nil
}
