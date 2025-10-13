package ports

import (
	"context"
	"go-dentist-api/internal/core/domain"
)

// PatientRepository define as operações de persistência para pacientes
type PatientRepository interface {
	// Operações básicas CRUD
	Create(ctx context.Context, patient *domain.Patient) error
	GetByID(ctx context.Context, id string) (*domain.Patient, error)
	GetByCPF(ctx context.Context, cpf string) (*domain.Patient, error)
	Update(ctx context.Context, patient *domain.Patient) error
	Delete(ctx context.Context, id string) error
	
	// Operações de listagem e busca
	GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error)
	GetActiveByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error)
	SearchByName(ctx context.Context, clinicID, name string) ([]*domain.Patient, error)
	
	// Operações de procedimentos
	AddProcedure(ctx context.Context, patientID string, procedure *domain.PatientProcedure) error
	GetProceduresByPatientID(ctx context.Context, patientID string) ([]*domain.PatientProcedure, error)
	GetProceduresByDentistID(ctx context.Context, dentistID string) ([]*domain.PatientProcedure, error)
	UpdateProcedure(ctx context.Context, procedure *domain.PatientProcedure) error
	
	// Operações de exames
	AddExam(ctx context.Context, patientID string, exam *domain.Exam) error
	GetExamsByPatientID(ctx context.Context, patientID string) ([]*domain.Exam, error)
	UpdateExam(ctx context.Context, exam *domain.Exam) error
	GetExamByID(ctx context.Context, examID string) (*domain.Exam, error)
	
	// Operações de histórico médico
	UpdateMedicalHistory(ctx context.Context, patientID string, medicalHistory *domain.MedicalHistory) error
	GetMedicalHistory(ctx context.Context, patientID string) (*domain.MedicalHistory, error)
}