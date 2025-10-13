package handlers

import (
	"encoding/json"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PatientHandler struct {
	patientService ports.PatientService
}

func NewPatientHandler(patientService ports.PatientService) *PatientHandler {
	return &PatientHandler{
		patientService: patientService,
	}
}

// CreatePatientRequest representa a requisição para criar um paciente
type CreatePatientRequest struct {
	ClinicID         string                   `json:"clinic_id"`
	Name             string                   `json:"name"`
	CPF              string                   `json:"cpf"`
	RG               string                   `json:"rg,omitempty"`
	Phone            string                   `json:"phone"`
	Email            string                   `json:"email,omitempty"`
	Address          domain.Address           `json:"address"`
	BirthDate        time.Time                `json:"birth_date"`
	Gender           string                   `json:"gender,omitempty"`
	MaritalStatus    string                   `json:"marital_status,omitempty"`
	Profession       string                   `json:"profession,omitempty"`
	EmergencyContact domain.EmergencyContact  `json:"emergency_contact"`
}

// UpdatePatientRequest representa a requisição para atualizar um paciente
type UpdatePatientRequest struct {
	Name             string                   `json:"name,omitempty"`
	Phone            string                   `json:"phone,omitempty"`
	Email            string                   `json:"email,omitempty"`
	Address          *domain.Address          `json:"address,omitempty"`
	Gender           string                   `json:"gender,omitempty"`
	MaritalStatus    string                   `json:"marital_status,omitempty"`
	Profession       string                   `json:"profession,omitempty"`
	EmergencyContact *domain.EmergencyContact `json:"emergency_contact,omitempty"`
	Active           *bool                    `json:"active,omitempty"`
}

// AddProcedureRequest representa a requisição para adicionar um procedimento
type AddProcedureRequest struct {
	DentistID     string  `json:"dentist_id"`
	DentistName   string  `json:"dentist_name"`
	ProcedureName string  `json:"procedure_name"`
	Description   string  `json:"description"`
	Cost          float64 `json:"cost"`
}

// AddExamRequest representa a requisição para adicionar um exame
type AddExamRequest struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	DentistID   string `json:"dentist_id"`
	DentistName string `json:"dentist_name"`
}

// UpdateExamRequest representa a requisição para atualizar um exame
type UpdateExamRequest struct {
	Results   string `json:"results,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	ReportURL string `json:"report_url,omitempty"`
	Status    string `json:"status,omitempty"`
}

// UpdateMedicalHistoryRequest representa a requisição para atualizar histórico médico
type UpdateMedicalHistoryRequest struct {
	Allergies         []string `json:"allergies,omitempty"`
	Medications       []string `json:"medications,omitempty"`
	MedicalConditions []string `json:"medical_conditions,omitempty"`
	ChronicDiseases   []string `json:"chronic_diseases,omitempty"`
	PreviousSurgeries []string `json:"previous_surgeries,omitempty"`
	FamilyHistory     []string `json:"family_history,omitempty"`
	DentalHistory     string   `json:"dental_history,omitempty"`
	Notes             string   `json:"notes,omitempty"`
	UpdatedBy         string   `json:"updated_by"`
}

// CreatePatient godoc
// @Summary Create a new patient
// @Description Create a new patient with the provided information
// @Tags patients
// @Accept json
// @Produce json
// @Param patient body CreatePatientRequest true "Patient data"
// @Success 201 {object} domain.Patient
// @Failure 400 {string} string "Invalid request body"
// @Router /patients [post]
func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
	var req CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	patient := &domain.Patient{
		ID:               uuid.New().String(),
		ClinicID:         req.ClinicID,
		Name:             req.Name,
		CPF:              req.CPF,
		RG:               req.RG,
		Phone:            req.Phone,
		Email:            req.Email,
		Address:          req.Address,
		BirthDate:        req.BirthDate,
		Gender:           req.Gender,
		MaritalStatus:    req.MaritalStatus,
		Profession:       req.Profession,
		EmergencyContact: req.EmergencyContact,
		Active:           true,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		MedicalHistory: domain.MedicalHistory{
			LastUpdate: time.Now(),
		},
	}

	if err := h.patientService.CreatePatient(r.Context(), patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// GetPatient godoc
// @Summary Get a patient by ID
// @Description Get patient details by ID
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} domain.Patient
// @Failure 404 {string} string "Patient not found"
// @Router /patients/{id} [get]
func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	patient, err := h.patientService.GetPatientByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// GetPatientByCPF godoc
// @Summary Get a patient by CPF
// @Description Get patient details by CPF
// @Tags patients
// @Accept json
// @Produce json
// @Param cpf path string true "Patient CPF"
// @Success 200 {object} domain.Patient
// @Failure 404 {string} string "Patient not found"
// @Router /patients/cpf/{cpf} [get]
func (h *PatientHandler) GetPatientByCPF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cpf := vars["cpf"]
	
	if cpf == "" {
		http.Error(w, "CPF parameter is required", http.StatusBadRequest)
		return
	}

	patient, err := h.patientService.GetPatientByCPF(r.Context(), cpf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// UpdatePatient godoc
// @Summary Update a patient
// @Description Update patient information
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param patient body UpdatePatientRequest true "Patient data"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request"
// @Router /patients/{id} [put]
func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Buscar o paciente existente
	patient, err := h.patientService.GetPatientByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if patient == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	// Atualizar campos fornecidos
	if req.Name != "" {
		patient.Name = req.Name
	}
	if req.Phone != "" {
		patient.Phone = req.Phone
	}
	if req.Email != "" {
		patient.Email = req.Email
	}
	if req.Address != nil {
		patient.Address = *req.Address
	}
	if req.Gender != "" {
		patient.Gender = req.Gender
	}
	if req.MaritalStatus != "" {
		patient.MaritalStatus = req.MaritalStatus
	}
	if req.Profession != "" {
		patient.Profession = req.Profession
	}
	if req.EmergencyContact != nil {
		patient.EmergencyContact = *req.EmergencyContact
	}
	if req.Active != nil {
		patient.Active = *req.Active
	}

	patient.UpdatedAt = time.Now()

	if err := h.patientService.UpdatePatient(r.Context(), patient); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeletePatient godoc
// @Summary Delete a patient
// @Description Delete a patient by ID
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request"
// @Router /patients/{id} [delete]
func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.patientService.DeletePatient(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeactivatePatient godoc
// @Summary Deactivate a patient
// @Description Deactivate a patient by ID
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request"
// @Router /patients/{id}/deactivate [patch]
func (h *PatientHandler) DeactivatePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.patientService.DeactivatePatient(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetPatientsByClinic godoc
// @Summary Get patients by clinic
// @Description Get all patients from a specific clinic
// @Tags patients
// @Accept json
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Param active query boolean false "Filter by active status"
// @Success 200 {array} domain.Patient
// @Failure 500 {string} string "Internal server error"
// @Router /clinics/{clinicId}/patients [get]
func (h *PatientHandler) GetPatientsByClinic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicID := vars["clinicId"]

	activeOnly := r.URL.Query().Get("active")
	var patients []*domain.Patient
	var err error

	if activeOnly == "true" {
		patients, err = h.patientService.GetActivePatientsByClinicID(r.Context(), clinicID)
	} else {
		patients, err = h.patientService.GetPatientsByClinicID(r.Context(), clinicID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// SearchPatients godoc
// @Summary Search patients by name
// @Description Search patients by name within a clinic
// @Tags patients
// @Accept json
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Param name query string true "Patient name to search"
// @Success 200 {array} domain.Patient
// @Failure 400 {string} string "Invalid request"
// @Router /clinics/{clinicId}/patients/search [get]
func (h *PatientHandler) SearchPatients(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicID := vars["clinicId"]
	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	patients, err := h.patientService.SearchPatientsByName(r.Context(), clinicID, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// AddProcedure godoc
// @Summary Add procedure to patient
// @Description Add a new procedure to patient's history
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param procedure body AddProcedureRequest true "Procedure data"
// @Success 201 "Created"
// @Failure 400 {string} string "Invalid request"
// @Router /patients/{id}/procedures [post]
func (h *PatientHandler) AddProcedure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	var req AddProcedureRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.patientService.AddProcedureToPatient(r.Context(), patientID, req.DentistID, req.DentistName, req.ProcedureName, req.Description, req.Cost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetPatientProcedures godoc
// @Summary Get patient procedures
// @Description Get all procedures for a specific patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {array} domain.PatientProcedure
// @Failure 404 {string} string "Patient not found"
// @Router /patients/{id}/procedures [get]
func (h *PatientHandler) GetPatientProcedures(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	procedures, err := h.patientService.GetPatientProcedures(r.Context(), patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procedures)
}

// AddExam godoc
// @Summary Add exam to patient
// @Description Add a new exam to patient's history
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param exam body AddExamRequest true "Exam data"
// @Success 201 "Created"
// @Failure 400 {string} string "Invalid request"
// @Router /patients/{id}/exams [post]
func (h *PatientHandler) AddExam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	var req AddExamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.patientService.AddExamToPatient(r.Context(), patientID, req.Type, req.Description, req.DentistID, req.DentistName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetPatientExams godoc
// @Summary Get patient exams
// @Description Get all exams for a specific patient
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {array} domain.Exam
// @Failure 404 {string} string "Patient not found"
// @Router /patients/{id}/exams [get]
func (h *PatientHandler) GetPatientExams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	exams, err := h.patientService.GetPatientExams(r.Context(), patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exams)
}

// UpdateExam godoc
// @Summary Update exam
// @Description Update exam results and status
// @Tags patients
// @Accept json
// @Produce json
// @Param examId path string true "Exam ID"
// @Param exam body UpdateExamRequest true "Exam data"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request"
// @Router /exams/{examId} [put]
func (h *PatientHandler) UpdateExam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	examID := vars["examId"]

	var req UpdateExamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "" {
		if err := h.patientService.UpdateExamStatus(r.Context(), examID, req.Status); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if req.Results != "" || req.ImageURL != "" || req.ReportURL != "" {
		if err := h.patientService.UpdateExamResults(r.Context(), examID, req.Results, req.ImageURL, req.ReportURL); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateMedicalHistory godoc
// @Summary Update patient medical history
// @Description Update patient's medical history
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param history body UpdateMedicalHistoryRequest true "Medical history data"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid request"
// @Router /patients/{id}/medical-history [put]
func (h *PatientHandler) UpdateMedicalHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	var req UpdateMedicalHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Combinar medical_conditions e chronic_diseases
	allConditions := append(req.MedicalConditions, req.ChronicDiseases...)
	
	// Converter FamilyHistory de []string para string
	familyHistoryStr := ""
	if len(req.FamilyHistory) > 0 {
		familyHistoryStr = strings.Join(req.FamilyHistory, "; ")
	}

	if err := h.patientService.UpdatePatientMedicalHistory(r.Context(), patientID, req.Allergies, req.Medications, allConditions, req.PreviousSurgeries, familyHistoryStr, req.DentalHistory, req.Notes, req.UpdatedBy); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMedicalHistory godoc
// @Summary Get patient medical history
// @Description Get patient's medical history
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} domain.MedicalHistory
// @Failure 404 {string} string "Patient not found"
// @Router /patients/{id}/medical-history [get]
func (h *PatientHandler) GetMedicalHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["id"]

	medicalHistory, err := h.patientService.GetPatientMedicalHistory(r.Context(), patientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(medicalHistory)
}

// GetPatientStatistics godoc
// @Summary Get patient statistics
// @Description Get statistics for patients in a clinic
// @Tags patients
// @Accept json
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {string} string "Internal server error"
// @Router /clinics/{clinicId}/patients/statistics [get]
func (h *PatientHandler) GetPatientStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicID := vars["clinicId"]

	stats, err := h.patientService.GetPatientStatistics(r.Context(), clinicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}