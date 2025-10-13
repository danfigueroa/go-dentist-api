package usecases

import (
	"context"
	"errors"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports/repositories"
	"go-dentist-api/internal/core/ports/services"
)

type ProcedureUseCase struct {
	procedureRepo repositories.ProcedureRepository
	clinicRepo    repositories.ClinicRepository
}

func NewProcedureUseCase(procedureRepo repositories.ProcedureRepository, clinicRepo repositories.ClinicRepository) services.ProcedureService {
	return &ProcedureUseCase{
		procedureRepo: procedureRepo,
		clinicRepo:    clinicRepo,
	}
}

func (uc *ProcedureUseCase) CreateProcedure(ctx context.Context, req services.CreateProcedureRequest) (*domain.Procedure, error) {
	clinic, err := uc.clinicRepo.GetByID(ctx, req.ClinicID)
	if err != nil {
		return nil, fmt.Errorf("clinic not found: %w", err)
	}
	if clinic == nil {
		return nil, errors.New("clinic not found")
	}

	procedure := domain.NewProcedure(req.ClinicID, req.AppointmentID, req.Name, req.Description, req.Price)

	if err := uc.procedureRepo.Create(ctx, procedure); err != nil {
		return nil, fmt.Errorf("failed to create procedure: %w", err)
	}

	return procedure, nil
}

func (uc *ProcedureUseCase) GetProcedure(ctx context.Context, id string) (*domain.Procedure, error) {
	procedure, err := uc.procedureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get procedure: %w", err)
	}

	return procedure, nil
}

func (uc *ProcedureUseCase) GetProceduresByAppointment(ctx context.Context, appointmentID string) ([]*domain.Procedure, error) {
	procedures, err := uc.procedureRepo.GetByAppointmentID(ctx, appointmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get procedures by appointment: %w", err)
	}

	return procedures, nil
}

func (uc *ProcedureUseCase) GetProceduresByClinic(ctx context.Context, clinicID string) ([]*domain.Procedure, error) {
	procedures, err := uc.procedureRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get procedures by clinic: %w", err)
	}

	return procedures, nil
}

func (uc *ProcedureUseCase) GetProceduresByName(ctx context.Context, name string) ([]*domain.Procedure, error) {
	procedures, err := uc.procedureRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get procedures by name: %w", err)
	}

	return procedures, nil
}

func (uc *ProcedureUseCase) UpdateProcedure(ctx context.Context, id string, req services.UpdateProcedureRequest) error {
	procedure, err := uc.procedureRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("procedure not found: %w", err)
	}

	if req.Name != "" {
		procedure.Name = req.Name
	}
	if req.Description != "" {
		procedure.Description = req.Description
	}
	if req.Price > 0 {
		procedure.Price = req.Price
	}

	if err := procedure.Validate(); err != nil {
		return fmt.Errorf("invalid procedure data: %w", err)
	}

	if err := uc.procedureRepo.Update(ctx, procedure); err != nil {
		return fmt.Errorf("failed to update procedure: %w", err)
	}

	return nil
}

func (uc *ProcedureUseCase) DeleteProcedure(ctx context.Context, id string) error {
	_, err := uc.procedureRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("procedure not found: %w", err)
	}

	if err := uc.procedureRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete procedure: %w", err)
	}

	return nil
}

func (uc *ProcedureUseCase) ListProcedures(ctx context.Context, filters map[string]interface{}) ([]*domain.Procedure, error) {
	procedures, err := uc.procedureRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list procedures: %w", err)
	}

	return procedures, nil
}
