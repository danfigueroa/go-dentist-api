package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Dentist struct {
	ID        string    `json:"id"`
	ClinicID  string    `json:"clinic_id"`
	Name      string    `json:"name"`
	CRO       string    `json:"cro"`
	Specialty string    `json:"specialty"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *Dentist) Validate() error {
	if d.Name == "" {
		return errors.New("dentist name is required")
	}
	if d.CRO == "" {
		return errors.New("CRO is required")
	}
	if d.Email == "" {
		return errors.New("email is required")
	}
	if d.ClinicID == "" {
		return errors.New("clinic ID is required")
	}
	return nil
}

func NewDentist(clinicID, name, cro, specialty, phone, email string) (*Dentist, error) {
	dentist := &Dentist{
		ID:        uuid.New().String(),
		ClinicID:  clinicID,
		Name:      name,
		CRO:       cro,
		Specialty: specialty,
		Phone:     phone,
		Email:     email,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := dentist.Validate(); err != nil {
		return nil, err
	}

	return dentist, nil
}
