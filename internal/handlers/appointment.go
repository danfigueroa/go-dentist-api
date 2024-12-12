package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-dentist-api/internal/config"
	"go-dentist-api/internal/models"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateAppointment godoc
// @Summary Create a new appointment
// @Description Create a new appointment by providing the details
// @Tags appointments
// @Accept json
// @Produce json
// @Param appointment body models.Appointment true "Appointment data"
// @Success 201 {object} models.Appointment
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 409 {string} string "Appointment with this ID already exists"
// @Failure 500 {string} string "Failed to save appointment"
// @Router /appointment [post]
func CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var appointment models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if appointment.ID == "" {
		appointment.ID = uuid.NewString()
	}

	if err := appointment.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if appointment.CreatedAt == "" {
		appointment.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if appointment.UpdatedAt == "" {
		appointment.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Appointments"),
		Item: map[string]types.AttributeValue{
			"ID":        &types.AttributeValueMemberS{Value: appointment.ID},
			"DentistID": &types.AttributeValueMemberS{Value: appointment.DentistID},
			"PatientID": &types.AttributeValueMemberS{Value: appointment.PatientID},
			"DateTime":  &types.AttributeValueMemberS{Value: appointment.DateTime},
			"Notes":     &types.AttributeValueMemberS{Value: appointment.Notes},
			"CreatedAt": &types.AttributeValueMemberS{Value: appointment.CreatedAt},
			"UpdatedAt": &types.AttributeValueMemberS{Value: appointment.UpdatedAt},
		},
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Appointment with this ID already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to save appointment", http.StatusInternalServerError)
		log.Printf("Error saving appointment: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(appointment)
}

// GetAllAppointments godoc
// @Summary Get all appointments
// @Description Retrieve all registered appointments
// @Tags appointments
// @Produce json
// @Success 200 {array} models.Appointment
// @Failure 500 {string} string "Failed to retrieve appointments"
// @Router /appointments [get]
func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Appointments"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve appointments", http.StatusInternalServerError)
		log.Printf("Error fetching appointments: %v", err)
		return
	}

	var appointments []models.Appointment
	err = attributevalue.UnmarshalListOfMaps(result.Items, &appointments)
	if err != nil {
		http.Error(w, "Failed to unmarshal appointment data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling appointment data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// GetAppointmentByID godoc
// @Summary Get appointment by ID
// @Description Retrieve a single appointment by providing its ID
// @Tags appointments
// @Produce json
// @Param id path string true "Appointment ID"
// @Success 200 {object} models.Appointment
// @Failure 404 {string} string "Appointment not found"
// @Failure 500 {string} string "Failed to retrieve appointment"
// @Router /appointment/{id} [get]
func GetAppointmentByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Appointments"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve appointment", http.StatusInternalServerError)
		log.Printf("Error fetching appointment with ID %s: %v", id, err)
		return
	}

	if result.Item == nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	var appointment models.Appointment
	err = attributevalue.UnmarshalMap(result.Item, &appointment)
	if err != nil {
		http.Error(w, "Failed to unmarshal appointment data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling appointment data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointment)
}

// UpdateAppointment godoc
// @Summary Update an existing appointment
// @Description Update fields of an existing appointment by providing its ID
// @Tags appointments
// @Accept json
// @Produce json
// @Param id path string true "Appointment ID"
// @Param appointment body models.Appointment true "Appointment data (ID will be ignored)"
// @Success 200 {object} models.Appointment
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 404 {string} string "Appointment not found"
// @Failure 500 {string} string "Failed to update appointment"
// @Router /appointment/{id} [put]
func UpdateAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Appointments"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve appointment", http.StatusInternalServerError)
		log.Printf("Error fetching appointment with ID %s: %v", id, err)
		return
	}
	if result.Item == nil {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	var currentAppointment models.Appointment
	if err = attributevalue.UnmarshalMap(result.Item, &currentAppointment); err != nil {
		http.Error(w, "Failed to unmarshal appointment data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling appointment data: %v", err)
		return
	}

	var updatedData models.Appointment
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedData.DentistID != "" {
		currentAppointment.DentistID = updatedData.DentistID
	}
	if updatedData.PatientID != "" {
		currentAppointment.PatientID = updatedData.PatientID
	}
	if updatedData.DateTime != "" {
		currentAppointment.DateTime = updatedData.DateTime
	}
	if updatedData.Notes != "" {
		currentAppointment.Notes = updatedData.Notes
	}

	// Valida campos obrigatórios após atualização
	if err := currentAppointment.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentAppointment.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	_, err = config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Appointments"),
		Item: map[string]types.AttributeValue{
			"ID":        &types.AttributeValueMemberS{Value: currentAppointment.ID},
			"DentistID": &types.AttributeValueMemberS{Value: currentAppointment.DentistID},
			"PatientID": &types.AttributeValueMemberS{Value: currentAppointment.PatientID},
			"DateTime":  &types.AttributeValueMemberS{Value: currentAppointment.DateTime},
			"Notes":     &types.AttributeValueMemberS{Value: currentAppointment.Notes},
			"CreatedAt": &types.AttributeValueMemberS{Value: currentAppointment.CreatedAt},
			"UpdatedAt": &types.AttributeValueMemberS{Value: currentAppointment.UpdatedAt},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Appointment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update appointment", http.StatusInternalServerError)
		log.Printf("Error updating appointment: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentAppointment)
}

// DeleteAppointment godoc
// @Summary Delete an appointment by ID
// @Description Delete a single appointment by providing its ID
// @Tags appointments
// @Produce json
// @Param id path string true "Appointment ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Appointment not found"
// @Failure 500 {string} string "Failed to delete appointment"
// @Router /appointment/{id} [delete]
func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := config.DBClient.DeleteItem(r.Context(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Appointments"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Appointment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete appointment", http.StatusInternalServerError)
		log.Printf("Error deleting appointment: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
