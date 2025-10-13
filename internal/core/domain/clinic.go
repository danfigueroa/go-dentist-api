package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Clinic struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (c *Clinic) Validate() error {
	if c.Name == "" {
		return errors.New("clinic name is required")
	}
	if c.CNPJ == "" {
		return errors.New("clinic CNPJ is required")
	}
	return nil
}

func NewClinic(name, cnpj, address, phone, email string) (*Clinic, error) {
	clinic := &Clinic{
		ID:        uuid.New().String(),
		Name:      name,
		CNPJ:      cnpj,
		Address:   address,
		Phone:     phone,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := clinic.Validate(); err != nil {
		return nil, err
	}

	return clinic, nil
}
