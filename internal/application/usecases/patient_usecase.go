package usecases

import (
	"context"
	"errors"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports"
)

type PatientUseCase struct {
	patientRepo ports.PatientRepository
}

func NewPatientUseCase(patientRepo ports.PatientRepository) ports.PatientService {
	return &PatientUseCase{
		patientRepo: patientRepo,
	}
}

func (uc *PatientUseCase) CreatePatient(ctx context.Context, patient *domain.Patient) error {
	if patient == nil {
		return errors.New("patient is required")
	}

	// Verificar se já existe um paciente com o mesmo CPF
	existingPatient, err := uc.patientRepo.GetByCPF(ctx, patient.CPF)
	if err != nil {
		return fmt.Errorf("failed to check existing patient: %w", err)
	}
	if existingPatient != nil {
		return errors.New("patient with this CPF already exists")
	}

	if err := patient.Validate(); err != nil {
		return fmt.Errorf("invalid patient data: %w", err)
	}

	if err := uc.patientRepo.Create(ctx, patient); err != nil {
		return fmt.Errorf("failed to create patient: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) GetPatientByID(ctx context.Context, id string) (*domain.Patient, error) {
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

func (uc *PatientUseCase) UpdatePatient(ctx context.Context, patient *domain.Patient) error {
	if patient == nil {
		return errors.New("patient is required")
	}
	if patient.ID == "" {
		return errors.New("patient ID is required")
	}

	// Verificar se o paciente existe
	existingPatient, err := uc.patientRepo.GetByID(ctx, patient.ID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if existingPatient == nil {
		return errors.New("patient not found")
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

	// Verificar se o paciente existe
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

func (uc *PatientUseCase) DeactivatePatient(ctx context.Context, id string) error {
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

	patient.Active = false
	if err := uc.patientRepo.Update(ctx, patient); err != nil {
		return fmt.Errorf("failed to deactivate patient: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) GetPatientsByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error) {
	if clinicID == "" {
		return nil, errors.New("clinic ID is required")
	}

	patients, err := uc.patientRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients by clinic: %w", err)
	}

	return patients, nil
}

func (uc *PatientUseCase) GetActivePatientsByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error) {
	if clinicID == "" {
		return nil, errors.New("clinic ID is required")
	}

	patients, err := uc.patientRepo.GetActiveByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active patients by clinic: %w", err)
	}

	return patients, nil
}

func (uc *PatientUseCase) SearchPatientsByName(ctx context.Context, clinicID, name string) ([]*domain.Patient, error) {
	if clinicID == "" {
		return nil, errors.New("clinic ID is required")
	}
	if name == "" {
		return nil, errors.New("name is required")
	}

	patients, err := uc.patientRepo.SearchByName(ctx, clinicID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to search patients by name: %w", err)
	}

	return patients, nil
}

func (uc *PatientUseCase) AddProcedureToPatient(ctx context.Context, patientID, dentistID, dentistName, procedureName, description string, cost float64) error {
	if patientID == "" {
		return errors.New("patient ID is required")
	}
	if dentistID == "" {
		return errors.New("dentist ID is required")
	}
	if procedureName == "" {
		return errors.New("procedure name is required")
	}

	procedure := &domain.PatientProcedure{
		DentistID:     dentistID,
		DentistName:   dentistName,
		ProcedureName: procedureName,
		Description:   description,
		Cost:          cost,
		Status:        "completed",
	}

	if err := uc.patientRepo.AddProcedure(ctx, patientID, procedure); err != nil {
		return fmt.Errorf("failed to add procedure to patient: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) GetPatientProcedures(ctx context.Context, patientID string) ([]*domain.PatientProcedure, error) {
	if patientID == "" {
		return nil, errors.New("patient ID is required")
	}

	procedures, err := uc.patientRepo.GetProceduresByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient procedures: %w", err)
	}

	return procedures, nil
}

func (uc *PatientUseCase) GetProceduresByDentist(ctx context.Context, dentistID string) ([]*domain.PatientProcedure, error) {
	if dentistID == "" {
		return nil, errors.New("dentist ID is required")
	}

	procedures, err := uc.patientRepo.GetProceduresByDentistID(ctx, dentistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get procedures by dentist: %w", err)
	}

	return procedures, nil
}

func (uc *PatientUseCase) UpdateProcedureStatus(ctx context.Context, procedureID, status string) error {
	if procedureID == "" {
		return errors.New("procedure ID is required")
	}
	if status == "" {
		return errors.New("status is required")
	}

	// Buscar o procedimento através de todos os pacientes
	// Isso é uma limitação do design atual, mas funciona
	return errors.New("procedure update not implemented - requires procedure ID lookup")
}

func (uc *PatientUseCase) AddExamToPatient(ctx context.Context, patientID, examType, description, dentistID, dentistName string) error {
	if patientID == "" {
		return errors.New("patient ID is required")
	}
	if examType == "" {
		return errors.New("exam type is required")
	}
	if dentistID == "" {
		return errors.New("dentist ID is required")
	}

	exam := &domain.Exam{
		Type:        examType,
		Description: description,
		DentistID:   dentistID,
		DentistName: dentistName,
		Status:      "pending",
	}

	if err := uc.patientRepo.AddExam(ctx, patientID, exam); err != nil {
		return fmt.Errorf("failed to add exam to patient: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) GetPatientExams(ctx context.Context, patientID string) ([]*domain.Exam, error) {
	if patientID == "" {
		return nil, errors.New("patient ID is required")
	}

	exams, err := uc.patientRepo.GetExamsByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient exams: %w", err)
	}

	return exams, nil
}

func (uc *PatientUseCase) UpdateExamResults(ctx context.Context, examID, results, imageURL, reportURL string) error {
	if examID == "" {
		return errors.New("exam ID is required")
	}

	exam, err := uc.patientRepo.GetExamByID(ctx, examID)
	if err != nil {
		return fmt.Errorf("failed to get exam: %w", err)
	}
	if exam == nil {
		return errors.New("exam not found")
	}

	exam.Results = results
	exam.ImageURL = imageURL
	exam.ReportURL = reportURL
	exam.Status = "completed"

	if err := uc.patientRepo.UpdateExam(ctx, exam); err != nil {
		return fmt.Errorf("failed to update exam results: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) UpdateExamStatus(ctx context.Context, examID, status string) error {
	if examID == "" {
		return errors.New("exam ID is required")
	}
	if status == "" {
		return errors.New("status is required")
	}

	exam, err := uc.patientRepo.GetExamByID(ctx, examID)
	if err != nil {
		return fmt.Errorf("failed to get exam: %w", err)
	}
	if exam == nil {
		return errors.New("exam not found")
	}

	exam.Status = status

	if err := uc.patientRepo.UpdateExam(ctx, exam); err != nil {
		return fmt.Errorf("failed to update exam status: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) UpdatePatientMedicalHistory(ctx context.Context, patientID string, allergies, medications, chronicDiseases, previousSurgeries []string, familyHistory, dentalHistory, notes, updatedBy string) error {
	if patientID == "" {
		return errors.New("patient ID is required")
	}

	medicalHistory := &domain.MedicalHistory{
		Allergies:         allergies,
		Medications:       medications,
		ChronicDiseases:   chronicDiseases,
		PreviousSurgeries: previousSurgeries,
		FamilyHistory:     familyHistory,
		DentalHistory:     dentalHistory,
		Notes:             notes,
		UpdatedBy:         updatedBy,
	}

	if err := uc.patientRepo.UpdateMedicalHistory(ctx, patientID, medicalHistory); err != nil {
		return fmt.Errorf("failed to update medical history: %w", err)
	}

	return nil
}

func (uc *PatientUseCase) GetPatientMedicalHistory(ctx context.Context, patientID string) (*domain.MedicalHistory, error) {
	if patientID == "" {
		return nil, errors.New("patient ID is required")
	}

	medicalHistory, err := uc.patientRepo.GetMedicalHistory(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get medical history: %w", err)
	}

	return medicalHistory, nil
}

func (uc *PatientUseCase) GetPatientStatistics(ctx context.Context, clinicID string) (map[string]interface{}, error) {
	if clinicID == "" {
		return nil, errors.New("clinic ID is required")
	}

	patients, err := uc.patientRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patients: %w", err)
	}

	activePatients := 0
	totalProcedures := 0
	totalExams := 0

	for _, patient := range patients {
		if patient.Active {
			activePatients++
		}
		totalProcedures += len(patient.Procedures)
		totalExams += len(patient.Exams)
	}

	stats := map[string]interface{}{
		"total_patients":    len(patients),
		"active_patients":   activePatients,
		"inactive_patients": len(patients) - activePatients,
		"total_procedures":  totalProcedures,
		"total_exams":       totalExams,
	}

	return stats, nil
}

func (uc *PatientUseCase) GetPatientProcedureHistory(ctx context.Context, patientID string) ([]*domain.PatientProcedure, error) {
	if patientID == "" {
		return nil, errors.New("patient ID is required")
	}

	procedures, err := uc.patientRepo.GetProceduresByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient procedure history: %w", err)
	}

	return procedures, nil
}
