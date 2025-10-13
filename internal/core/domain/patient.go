package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Address representa o endereço do paciente
type Address struct {
	Street     string `json:"street"`
	Number     string `json:"number"`
	Complement string `json:"complement,omitempty"`
	District   string `json:"district"`
	City       string `json:"city"`
	State      string `json:"state"`
	ZipCode    string `json:"zip_code"`
	Country    string `json:"country"`
}

// MedicalHistory representa o histórico médico do paciente
type MedicalHistory struct {
	Allergies           []string  `json:"allergies,omitempty"`
	Medications         []string  `json:"medications,omitempty"`
	ChronicDiseases     []string  `json:"chronic_diseases,omitempty"`
	PreviousSurgeries   []string  `json:"previous_surgeries,omitempty"`
	FamilyHistory       string    `json:"family_history,omitempty"`
	DentalHistory       string    `json:"dental_history,omitempty"`
	Notes               string    `json:"notes,omitempty"`
	LastUpdate          time.Time `json:"last_update"`
	UpdatedBy           string    `json:"updated_by"`
}

// PatientProcedure representa um procedimento realizado no paciente
type PatientProcedure struct {
	ID            string    `json:"id"`
	PatientID     string    `json:"patient_id"`
	DentistID     string    `json:"dentist_id"`
	DentistName   string    `json:"dentist_name"`
	ProcedureName string    `json:"procedure_name"`
	Description   string    `json:"description"`
	Date          time.Time `json:"date"`
	Status        string    `json:"status"` // completed, in_progress, cancelled
	Cost          float64   `json:"cost"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// Exam representa exames realizados pelo paciente
type Exam struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patient_id"`
	Type        string    `json:"type"` // radiografia, tomografia, etc.
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Results     string    `json:"results,omitempty"`
	ImageURL    string    `json:"image_url,omitempty"`
	ReportURL   string    `json:"report_url,omitempty"`
	DentistID   string    `json:"dentist_id"`
	DentistName string    `json:"dentist_name"`
	Status      string    `json:"status"` // pending, completed, reviewed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EmergencyContact representa contato de emergência do paciente
type EmergencyContact struct {
	Name         string `json:"name"`
	Relationship string `json:"relationship"`
	Phone        string `json:"phone"`
	Email        string `json:"email,omitempty"`
}

// Patient representa um paciente da clínica
type Patient struct {
	ID               string            `json:"id"`
	ClinicID         string            `json:"clinic_id"`
	Name             string            `json:"name"`
	CPF              string            `json:"cpf"`
	RG               string            `json:"rg,omitempty"`
	Phone            string            `json:"phone"`
	Email            string            `json:"email,omitempty"`
	Address          Address           `json:"address"`
	BirthDate        time.Time         `json:"birth_date"`
	Gender           string            `json:"gender,omitempty"` // M, F, Other
	MaritalStatus    string            `json:"marital_status,omitempty"`
	Profession       string            `json:"profession,omitempty"`
	EmergencyContact EmergencyContact  `json:"emergency_contact"`
	MedicalHistory   MedicalHistory    `json:"medical_history"`
	Procedures       []PatientProcedure `json:"procedures,omitempty"`
	Exams            []Exam            `json:"exams,omitempty"`
	Active           bool              `json:"active"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
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
	if p.Phone == "" {
		return errors.New("phone is required")
	}
	if p.Address.Street == "" || p.Address.City == "" || p.Address.State == "" {
		return errors.New("complete address is required")
	}
	if p.EmergencyContact.Name == "" || p.EmergencyContact.Phone == "" {
		return errors.New("emergency contact information is required")
	}
	return nil
}

func NewPatient(clinicID, name, cpf, phone string) (*Patient, error) {
	patient := &Patient{
		ID:        uuid.New().String(),
		ClinicID:  clinicID,
		Name:      name,
		CPF:       cpf,
		Phone:     phone,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		MedicalHistory: MedicalHistory{
			LastUpdate: time.Now(),
		},
	}

	if err := patient.Validate(); err != nil {
		return nil, err
	}

	return patient, nil
}

// AddProcedure adiciona um procedimento ao histórico do paciente
func (p *Patient) AddProcedure(dentistID, dentistName, procedureName, description string, cost float64) {
	procedure := PatientProcedure{
		ID:            uuid.New().String(),
		PatientID:     p.ID,
		DentistID:     dentistID,
		DentistName:   dentistName,
		ProcedureName: procedureName,
		Description:   description,
		Date:          time.Now(),
		Status:        "completed",
		Cost:          cost,
		CreatedAt:     time.Now(),
	}
	p.Procedures = append(p.Procedures, procedure)
	p.UpdatedAt = time.Now()
}

// AddExam adiciona um exame ao histórico do paciente
func (p *Patient) AddExam(examType, description, dentistID, dentistName string) {
	exam := Exam{
		ID:          uuid.New().String(),
		PatientID:   p.ID,
		Type:        examType,
		Description: description,
		Date:        time.Now(),
		DentistID:   dentistID,
		DentistName: dentistName,
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	p.Exams = append(p.Exams, exam)
	p.UpdatedAt = time.Now()
}

// UpdateMedicalHistory atualiza o histórico médico do paciente
func (p *Patient) UpdateMedicalHistory(allergies, medications, chronicDiseases, previousSurgeries []string, familyHistory, dentalHistory, notes, updatedBy string) {
	p.MedicalHistory = MedicalHistory{
		Allergies:         allergies,
		Medications:       medications,
		ChronicDiseases:   chronicDiseases,
		PreviousSurgeries: previousSurgeries,
		FamilyHistory:     familyHistory,
		DentalHistory:     dentalHistory,
		Notes:             notes,
		LastUpdate:        time.Now(),
		UpdatedBy:         updatedBy,
	}
	p.UpdatedAt = time.Now()
}
