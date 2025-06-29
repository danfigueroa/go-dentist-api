package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"dental-saas/modules/dental/models"
	"dental-saas/shared/config"
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
// @Router /api/v1/dental/appointment [post]
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

	item := map[string]types.AttributeValue{
		"ID":        &types.AttributeValueMemberS{Value: appointment.ID},
		"PatientID": &types.AttributeValueMemberS{Value: appointment.PatientID},
		"DentistID": &types.AttributeValueMemberS{Value: appointment.DentistID},
		"DateTime":  &types.AttributeValueMemberS{Value: appointment.DateTime},
		"Status":    &types.AttributeValueMemberS{Value: appointment.Status},
		"CreatedAt": &types.AttributeValueMemberS{Value: appointment.CreatedAt},
		"UpdatedAt": &types.AttributeValueMemberS{Value: appointment.UpdatedAt},
	}

	if appointment.ProcedureID != "" {
		item["ProcedureID"] = &types.AttributeValueMemberS{Value: appointment.ProcedureID}
	}
	if appointment.Notes != "" {
		item["Notes"] = &types.AttributeValueMemberS{Value: appointment.Notes}
	}
	if appointment.Duration != "" {
		item["Duration"] = &types.AttributeValueMemberS{Value: appointment.Duration}
	}

	_, err := config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName:           aws.String("Appointments"),
		Item:                item,
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
// @Description Get a list of all appointments
// @Tags appointments
// @Produce json
// @Success 200 {array} models.Appointment
// @Failure 500 {string} string "Failed to retrieve appointments"
// @Router /api/v1/dental/appointment [get]
func GetAllAppointments(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Appointments"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve appointments", http.StatusInternalServerError)
		log.Printf("Error scanning appointments: %v", err)
		return
	}

	var appointments []models.Appointment
	for _, item := range result.Items {
		var appointment models.Appointment
		if err := attributevalue.UnmarshalMap(item, &appointment); err != nil {
			log.Printf("Error unmarshaling appointment: %v", err)
			continue
		}
		appointments = append(appointments, appointment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// GetAppointmentByID godoc
// @Summary Get appointment by ID
// @Description Get an appointment by its ID
// @Tags appointments
// @Produce json
// @Param id path string true "Appointment ID"
// @Success 200 {object} models.Appointment
// @Failure 404 {string} string "Appointment not found"
// @Failure 500 {string} string "Failed to retrieve appointment"
// @Router /api/v1/dental/appointment/{id} [get]
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
	if err = attributevalue.UnmarshalMap(result.Item, &appointment); err != nil {
		http.Error(w, "Failed to unmarshal appointment data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling appointment data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointment)
}

// GetAppointmentsByPatient godoc
// @Summary Get appointments by patient ID
// @Description Get all appointments for a specific patient
// @Tags appointments
// @Produce json
// @Param patientId path string true "Patient ID"
// @Success 200 {array} models.Appointment
// @Failure 500 {string} string "Failed to retrieve appointments"
// @Router /api/v1/dental/appointment/patient/{patientId} [get]
func GetAppointmentsByPatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	patientID := vars["patientId"]

	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("Appointments"),
		FilterExpression: aws.String("PatientID = :patientId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":patientId": &types.AttributeValueMemberS{Value: patientID},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve appointments", http.StatusInternalServerError)
		log.Printf("Error scanning appointments by patient: %v", err)
		return
	}

	var appointments []models.Appointment
	for _, item := range result.Items {
		var appointment models.Appointment
		if err := attributevalue.UnmarshalMap(item, &appointment); err != nil {
			log.Printf("Error unmarshaling appointment: %v", err)
			continue
		}
		appointments = append(appointments, appointment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
}

// GetAppointmentsByDentist godoc
// @Summary Get appointments by dentist ID
// @Description Get all appointments for a specific dentist
// @Tags appointments
// @Produce json
// @Param dentistId path string true "Dentist ID"
// @Success 200 {array} models.Appointment
// @Failure 500 {string} string "Failed to retrieve appointments"
// @Router /api/v1/dental/appointment/dentist/{dentistId} [get]
func GetAppointmentsByDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dentistID := vars["dentistId"]

	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("Appointments"),
		FilterExpression: aws.String("DentistID = :dentistId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dentistId": &types.AttributeValueMemberS{Value: dentistID},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve appointments", http.StatusInternalServerError)
		log.Printf("Error scanning appointments by dentist: %v", err)
		return
	}

	var appointments []models.Appointment
	for _, item := range result.Items {
		var appointment models.Appointment
		if err := attributevalue.UnmarshalMap(item, &appointment); err != nil {
			log.Printf("Error unmarshaling appointment: %v", err)
			continue
		}
		appointments = append(appointments, appointment)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(appointments)
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
// @Router /api/v1/dental/appointment/{id} [put]
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

	if updatedData.PatientID != "" {
		currentAppointment.PatientID = updatedData.PatientID
	}
	if updatedData.DentistID != "" {
		currentAppointment.DentistID = updatedData.DentistID
	}
	if updatedData.ProcedureID != "" {
		currentAppointment.ProcedureID = updatedData.ProcedureID
	}
	if updatedData.DateTime != "" {
		currentAppointment.DateTime = updatedData.DateTime
	}
	if updatedData.Duration != "" {
		currentAppointment.Duration = updatedData.Duration
	}
	if updatedData.Status != "" {
		currentAppointment.Status = updatedData.Status
	}
	if updatedData.Notes != "" {
		currentAppointment.Notes = updatedData.Notes
	}

	if err := currentAppointment.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentAppointment.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	item := map[string]types.AttributeValue{
		"ID":        &types.AttributeValueMemberS{Value: currentAppointment.ID},
		"PatientID": &types.AttributeValueMemberS{Value: currentAppointment.PatientID},
		"DentistID": &types.AttributeValueMemberS{Value: currentAppointment.DentistID},
		"DateTime":  &types.AttributeValueMemberS{Value: currentAppointment.DateTime},
		"Status":    &types.AttributeValueMemberS{Value: currentAppointment.Status},
		"CreatedAt": &types.AttributeValueMemberS{Value: currentAppointment.CreatedAt},
		"UpdatedAt": &types.AttributeValueMemberS{Value: currentAppointment.UpdatedAt},
	}

	if currentAppointment.ProcedureID != "" {
		item["ProcedureID"] = &types.AttributeValueMemberS{Value: currentAppointment.ProcedureID}
	}
	if currentAppointment.Notes != "" {
		item["Notes"] = &types.AttributeValueMemberS{Value: currentAppointment.Notes}
	}
	if currentAppointment.Duration != "" {
		item["Duration"] = &types.AttributeValueMemberS{Value: currentAppointment.Duration}
	}

	_, err = config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName:           aws.String("Appointments"),
		Item:                item,
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
// @Summary Delete an appointment
// @Description Delete an appointment by its ID
// @Tags appointments
// @Param id path string true "Appointment ID"
// @Success 204 "Appointment deleted successfully"
// @Failure 404 {string} string "Appointment not found"
// @Failure 500 {string} string "Failed to delete appointment"
// @Router /api/v1/dental/appointment/{id} [delete]
func DeleteAppointment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := config.DBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
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