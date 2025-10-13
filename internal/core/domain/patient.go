package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Patient struct {
	ID             string    `json:"id"`
	ClinicID       string    `json:"clinic_id"`
	Name           string    `json:"name"`
	CPF            string    `json:"cpf"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Address        string    `json:"address"`
	BirthDate      time.Time `json:"birth_date"`
	MedicalHistory string    `json:"medical_history"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (p *Patient) Validate() error {
	if p.Name == "" {
		return errors.New("patient name is required")
	}
	if p.CPF == "" {
		return errors.New("CPF is required")
	}
	if p.ClinicID == "" {
		return errors.New("clinic ID is required")
	}
	return nil
}

func NewPatient(clinicID, name, cpf, phone, email, address string, birthDate time.Time, medicalHistory string) (*Patient, error) {
	patient := &Patient{
		ID:             uuid.New().String(),
		ClinicID:       clinicID,
		Name:           name,
		CPF:            cpf,
		Phone:          phone,
		Email:          email,
		Address:        address,
		BirthDate:      birthDate,
		MedicalHistory: medicalHistory,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := patient.Validate(); err != nil {
		return nil, err
	}

	return patient, nil
}
