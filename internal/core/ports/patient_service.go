package ports

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

// PatientService define os casos de uso para pacientes
type PatientService interface {
	// Operações básicas de pacientes
	CreatePatient(ctx context.Context, patient *domain.Patient) error
	GetPatientByID(ctx context.Context, id string) (*domain.Patient, error)
	GetPatientByCPF(ctx context.Context, cpf string) (*domain.Patient, error)
	UpdatePatient(ctx context.Context, patient *domain.Patient) error
	DeletePatient(ctx context.Context, id string) error
	DeactivatePatient(ctx context.Context, id string) error
	
	// Operações de listagem e busca
	GetPatientsByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error)
	GetActivePatientsByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error)
	SearchPatientsByName(ctx context.Context, clinicID, name string) ([]*domain.Patient, error)
	
	// Operações de procedimentos
	AddProcedureToPatient(ctx context.Context, patientID, dentistID, dentistName, procedureName, description string, cost float64) error
	GetPatientProcedures(ctx context.Context, patientID string) ([]*domain.PatientProcedure, error)
	GetProceduresByDentist(ctx context.Context, dentistID string) ([]*domain.PatientProcedure, error)
	UpdateProcedureStatus(ctx context.Context, procedureID, status string) error
	
	// Operações de exames
	AddExamToPatient(ctx context.Context, patientID, examType, description, dentistID, dentistName string) error
	GetPatientExams(ctx context.Context, patientID string) ([]*domain.Exam, error)
	UpdateExamResults(ctx context.Context, examID, results, imageURL, reportURL string) error
	UpdateExamStatus(ctx context.Context, examID, status string) error
	
	// Operações de histórico médico
	UpdatePatientMedicalHistory(ctx context.Context, patientID string, allergies, medications, chronicDiseases, previousSurgeries []string, familyHistory, dentalHistory, notes, updatedBy string) error
	GetPatientMedicalHistory(ctx context.Context, patientID string) (*domain.MedicalHistory, error)
	
	// Relatórios e estatísticas
	GetPatientStatistics(ctx context.Context, clinicID string) (map[string]interface{}, error)
	GetPatientProcedureHistory(ctx context.Context, patientID string) ([]*domain.PatientProcedure, error)
}