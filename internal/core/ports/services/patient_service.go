package services

import (
	"context"
	"go-dentist-api/internal/core/domain"
	"time"
)

type CreatePatientRequest struct {
	ClinicID       string    `json:"clinic_id"`
	Name           string    `json:"name"`
	CPF            string    `json:"cpf"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	BirthDate      time.Time `json:"birth_date"`
	MedicalHistory string    `json:"medical_history"`
}

type UpdatePatientRequest struct {
	Name           string    `json:"name"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	BirthDate      time.Time `json:"birth_date"`
	MedicalHistory string    `json:"medical_history"`
}

type PatientService interface {
	CreatePatient(ctx context.Context, req CreatePatientRequest) (*domain.Patient, error)
	GetPatient(ctx context.Context, id string) (*domain.Patient, error)
	GetPatientByCPF(ctx context.Context, cpf string) (*domain.Patient, error)
	GetPatientsByName(ctx context.Context, name string) ([]*domain.Patient, error)
	GetPatientsByClinic(ctx context.Context, clinicID string) ([]*domain.Patient, error)
	UpdatePatient(ctx context.Context, id string, req UpdatePatientRequest) error
	DeletePatient(ctx context.Context, id string) error
	ListPatients(ctx context.Context, filters map[string]interface{}) ([]*domain.Patient, error)
}
