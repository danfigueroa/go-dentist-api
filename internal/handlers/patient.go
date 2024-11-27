package handlers

import (
	"encoding/json"
	"go-dentist-api/internal/models"
	"net/http"
)

func CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lógica para salvar o paciente no DynamoDB pode ser adicionada aqui.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}
