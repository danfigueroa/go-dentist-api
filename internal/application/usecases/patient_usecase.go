package usecases

import (
	"context"
	"errors"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports/repositories"
	"go-dentist-api/internal/core/ports/services"
)

type PatientUseCase struct {
	patientRepo repositories.PatientRepository
	clinicRepo  repositories.ClinicRepository
}

func NewPatientUseCase(patientRepo repositories.PatientRepository, clinicRepo repositories.ClinicRepository) services.PatientService {
	return &PatientUseCase{
		patientRepo: patientRepo,
		clinicRepo:  clinicRepo,
	}
}

func (uc *PatientUseCase) CreatePatient(ctx context.Context, req services.CreatePatientRequest) (*domain.Patient, error) {
	clinic, err := uc.clinicRepo.GetByID(ctx, req.ClinicID)
	if err != nil {
		return nil, fmt.Errorf("clinic not found: %w", err)
	}
	if clinic == nil {
		return nil, errors.New("clinic not found")
	}

	existingPatient, _ := uc.patientRepo.GetByCPF(ctx, req.CPF)
	if existingPatient != nil {
		return nil, errors.New("patient with this CPF already exists")
	}

	patient, err := domain.NewPatient(req.ClinicID, req.Name, req.CPF, req.Phone, req.Email, req.Address, req.BirthDate, req.MedicalHistory)
	if err != nil {
		return nil, fmt.Errorf("invalid patient data: %w", err)
	}

	if err := uc.patientRepo.Create(ctx, patient); err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	return patient, nil
}

func (uc *PatientUseCase) GetPatient(ctx context.Context, id string) (*domain.Patient, error) {
	if id == "" {
		return nil, errors.New("patient ID is required")
	}

	patient, err := uc.patientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	return patient, nil
}

func (uc *PatientUseCase) GetPatientByCPF(ctx context.Context, cpf string) (*domain.Patient, error) {
	if cpf == "" {
		return nil, errors.New("CPF is required")
	}

	patient, err := uc.patientRepo.GetByCPF(ctx, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient by CPF: %w", err)
	}

	return patient, nil
}

func (uc *PatientUseCase) GetPatientsByName(ctx context.Context, name string) ([]*domain.Patient, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	patients, err := uc.patientRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients by name: %w", err)
	}

	return patients, nil
}

func (uc *PatientUseCase) GetPatientsByClinic(ctx context.Context, clinicID string) ([]*domain.Patient, error) {
	if clinicID == "" {
		return nil, errors.New("clinic ID is required")
	}

	patients, err := uc.patientRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients by clinic: %w", err)
	}

	return patients, nil
}

func (uc *PatientUseCase) UpdatePatient(ctx context.Context, id string, req services.UpdatePatientRequest) error {
	if id == "" {
		return errors.New("patient ID is required")
	}

	patient, err := uc.patientRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return errors.New("patient not found")
	}

	if req.Name != "" {
		patient.Name = req.Name
	}
	if req.Phone != "" {
		patient.Phone = req.Phone
	}
	if req.Email != "" {
		patient.Email = req.Email
	}
	if req.Address != "" {
		patient.Address = req.Address
	}
	if !req.BirthDate.IsZero() {
		patient.BirthDate = req.BirthDate
	}
	if req.MedicalHistory != "" {
		patient.MedicalHistory = req.MedicalHistory
	}

	if err := patient.Validate(); err != nil {
		return fmt.Errorf("invalid patient data: %w", err)
	}

	if err := uc.patientRepo.Update(ctx, patient); err != nil {
		return fmt.Errorf("failed to update patient: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) DeletePatient(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("patient ID is required")
	}

	patient, err := uc.patientRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return errors.New("patient not found")
	}

	if err := uc.patientRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete patient: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) ListPatients(ctx context.Context, filters map[string]interface{}) ([]*domain.Patient, error) {
	patients, err := uc.patientRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list patients: %w", err)
	}

	return patients, nil
}
