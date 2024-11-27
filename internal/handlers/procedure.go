package handlers

import (
	"encoding/json"
	"go-dentist-api/internal/models"
	"net/http"
)

func CreateProcedure(w http.ResponseWriter, r *http.Request) {
	var procedure models.Procedure
	if err := json.NewDecoder(r.Body).Decode(&procedure); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lógica para salvar o procedimento no DynamoDB pode ser adicionada aqui.

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(procedure)
}
