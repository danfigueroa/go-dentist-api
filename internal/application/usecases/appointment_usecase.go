package usecases

import (
	"context"
	"errors"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports/repositories"
	"go-dentist-api/internal/core/ports/services"
	"time"
)

type AppointmentUseCase struct {
	appointmentRepo repositories.AppointmentRepository
	clinicRepo      repositories.ClinicRepository
	patientRepo     repositories.PatientRepository
	dentistRepo     repositories.DentistRepository
}

func NewAppointmentUseCase(
	appointmentRepo repositories.AppointmentRepository,
	clinicRepo repositories.ClinicRepository,
	patientRepo repositories.PatientRepository,
	dentistRepo repositories.DentistRepository,
) services.AppointmentService {
	return &AppointmentUseCase{
		appointmentRepo: appointmentRepo,
		clinicRepo:      clinicRepo,
		patientRepo:     patientRepo,
		dentistRepo:     dentistRepo,
	}
}

func (uc *AppointmentUseCase) CreateAppointment(ctx context.Context, req services.CreateAppointmentRequest) (*domain.Appointment, error) {
	clinic, err := uc.clinicRepo.GetByID(ctx, req.ClinicID)
	if err != nil {
		return nil, fmt.Errorf("clinic not found: %w", err)
	}
	if clinic == nil {
		return nil, errors.New("clinic not found")
	}

	patient, err := uc.patientRepo.GetByID(ctx, req.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient not found: %w", err)
	}
	if patient == nil {
		return nil, errors.New("patient not found")
	}

	dentist, err := uc.dentistRepo.GetByID(ctx, req.DentistID)
	if err != nil {
		return nil, fmt.Errorf("dentist not found: %w", err)
	}
	if dentist == nil {
		return nil, errors.New("dentist not found")
	}

	appointment := domain.NewAppointment(req.ClinicID, req.PatientID, req.DentistID, req.UserID, req.ScheduledAt, req.Notes)

	if err := uc.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	return appointment, nil
}

func (uc *AppointmentUseCase) GetAppointment(ctx context.Context, id string) (*domain.Appointment, error) {
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return appointment, nil
}

func (uc *AppointmentUseCase) GetAppointmentsByPatient(ctx context.Context, patientID string) ([]*domain.Appointment, error) {
	appointments, err := uc.appointmentRepo.GetByPatientID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by patient: %w", err)
	}

	return appointments, nil
}

func (uc *AppointmentUseCase) GetAppointmentsByDentist(ctx context.Context, dentistID string) ([]*domain.Appointment, error) {
	appointments, err := uc.appointmentRepo.GetByDentistID(ctx, dentistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by dentist: %w", err)
	}

	return appointments, nil
}

func (uc *AppointmentUseCase) GetAppointmentsByClinic(ctx context.Context, clinicID string) ([]*domain.Appointment, error) {
	appointments, err := uc.appointmentRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by clinic: %w", err)
	}

	return appointments, nil
}

func (uc *AppointmentUseCase) GetAppointmentsByDateRange(ctx context.Context, clinicID string, startDate, endDate time.Time) ([]*domain.Appointment, error) {
	appointments, err := uc.appointmentRepo.GetByDateRange(ctx, clinicID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by date range: %w", err)
	}

	return appointments, nil
}

func (uc *AppointmentUseCase) UpdateAppointment(ctx context.Context, id string, req services.UpdateAppointmentRequest) error {
	appointment, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("appointment not found: %w", err)
	}

	if !req.ScheduledAt.IsZero() {
		appointment.ScheduledAt = req.ScheduledAt
	}
	if req.Notes != "" {
		appointment.Notes = req.Notes
	}
	if req.Status != "" {
		appointment.Status = domain.AppointmentStatus(req.Status)
	}
	if req.TotalAmount > 0 {
		appointment.TotalAmount = req.TotalAmount
	}

	if err := appointment.Validate(); err != nil {
		return fmt.Errorf("invalid appointment data: %w", err)
	}

	if err := uc.appointmentRepo.Update(ctx, appointment); err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	return nil
}

func (uc *AppointmentUseCase) DeleteAppointment(ctx context.Context, id string) error {
	_, err := uc.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("appointment not found: %w", err)
	}

	if err := uc.appointmentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	return nil
}

func (uc *AppointmentUseCase) ListAppointments(ctx context.Context, filters map[string]interface{}) ([]*domain.Appointment, error) {
	appointments, err := uc.appointmentRepo.List(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appointments, nil
}
