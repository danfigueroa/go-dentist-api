package models

type Patient struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	DateOfBirth  string `json:"date_of_birth"`
	MedicalNotes string `json:"medical_notes"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
