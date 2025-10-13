package handlers

import (
	"encoding/json"
	"go-dentist-api/internal/core/ports/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DentistHandler struct {
	dentistService services.DentistService
}

func NewDentistHandler(dentistService services.DentistService) *DentistHandler {
	return &DentistHandler{
		dentistService: dentistService,
	}
}

// CreateDentist godoc
// @Summary Create a new dentist
// @Description Create a new dentist with the provided information
// @Tags dentists
// @Accept json
// @Produce json
// @Param dentist body services.CreateDentistRequest true "Dentist data"
// @Success 201 {object} domain.Dentist
// @Failure 400 {string} string "Invalid request body"
// @Router /dentists [post]
func (h *DentistHandler) CreateDentist(w http.ResponseWriter, r *http.Request) {
	var req services.CreateDentistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dentist, err := h.dentistService.CreateDentist(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dentist)
}

func (h *DentistHandler) GetDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	dentist, err := h.dentistService.GetDentist(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if dentist == nil {
		http.Error(w, "Dentist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentist)
}

func (h *DentistHandler) GetDentistByCRO(w http.ResponseWriter, r *http.Request) {
	cro := r.URL.Query().Get("cro")
	if cro == "" {
		http.Error(w, "CRO parameter is required", http.StatusBadRequest)
		return
	}

	dentist, err := h.dentistService.GetDentistByCRO(r.Context(), cro)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if dentist == nil {
		http.Error(w, "Dentist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentist)
}

func (h *DentistHandler) GetDentistsByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	dentists, err := h.dentistService.GetDentistsByName(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}

func (h *DentistHandler) GetDentistsByClinic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clinicID := vars["clinicId"]

	dentists, err := h.dentistService.GetDentistsByClinic(r.Context(), clinicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}

func (h *DentistHandler) UpdateDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req services.UpdateDentistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.dentistService.UpdateDentist(r.Context(), id, req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DentistHandler) DeleteDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.dentistService.DeleteDentist(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListDentists godoc
// @Summary List all dentists
// @Description Get a list of all dentists with optional filters
// @Tags dentists
// @Accept json
// @Produce json
// @Param clinic_id query string false "Filter by clinic ID"
// @Param active query boolean false "Filter by active status"
// @Param specialty query string false "Filter by specialty"
// @Success 200 {array} domain.Dentist
// @Failure 500 {string} string "Internal server error"
// @Router /dentists [get]
func (h *DentistHandler) ListDentists(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]interface{})

	if clinicID := r.URL.Query().Get("clinic_id"); clinicID != "" {
		filters["clinic_id"] = clinicID
	}

	if active := r.URL.Query().Get("active"); active != "" {
		if activeBool, err := strconv.ParseBool(active); err == nil {
			filters["active"] = activeBool
		}
	}

	if specialty := r.URL.Query().Get("specialty"); specialty != "" {
		filters["specialty"] = specialty
	}

	dentists, err := h.dentistService.ListDentists(r.Context(), filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}
