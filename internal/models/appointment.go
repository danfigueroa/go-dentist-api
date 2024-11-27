package models

type Appointment struct {
	ID        string `json:"id"`
	DentistID string `json:"dentist_id"`
	PatientID string `json:"patient_id"`
	DateTime  string `json:"date_time"`
	Notes     string `json:"notes"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
