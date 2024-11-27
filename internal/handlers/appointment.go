package handlers

import (
	"encoding/json"
	"go-dentist-api/internal/models"
	"net/http"
)

func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Aqui você adicionaria lógica para salvar o appointment no DynamoDB.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appointment)
}
