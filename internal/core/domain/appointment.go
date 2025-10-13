package domain

import (
	"errors"
	"time"
)

type AppointmentStatus string

const (
	AppointmentStatusScheduled AppointmentStatus = "scheduled"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCompleted AppointmentStatus = "completed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	ID          string            `json:"id"`
	ClinicID    string            `json:"clinic_id"`
	PatientID   string            `json:"patient_id"`
	DentistID   string            `json:"dentist_id"`
	UserID      string            `json:"user_id"`
	ScheduledAt time.Time         `json:"scheduled_at"`
	Status      AppointmentStatus `json:"status"`
	Notes       string            `json:"notes"`
	TotalAmount float64           `json:"total_amount"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (a *Appointment) Validate() error {
	if a.ClinicID == "" {
		return errors.New("clinic ID is required")
	}
	if a.PatientID == "" {
		return errors.New("patient ID is required")
	}
	if a.DentistID == "" {
		return errors.New("dentist ID is required")
	}
	if a.ScheduledAt.IsZero() {
		return errors.New("scheduled time is required")
	}
	return nil
}

func NewAppointment(clinicID, patientID, dentistID, userID string, scheduledAt time.Time, notes string) *Appointment {
	now := time.Now()
	return &Appointment{
		ClinicID:    clinicID,
		PatientID:   patientID,
		DentistID:   dentistID,
		UserID:      userID,
		ScheduledAt: scheduledAt,
		Status:      AppointmentStatusScheduled,
		Notes:       notes,
		TotalAmount: 0.0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
