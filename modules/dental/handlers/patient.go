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

// CreatePatient godoc
// @Summary Create a new patient
// @Description Create a new patient by providing the details
// @Tags patients
// @Accept json
// @Produce json
// @Param patient body models.Patient true "Patient data"
// @Success 201 {object} models.Patient
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 409 {string} string "Patient with this ID already exists"
// @Failure 500 {string} string "Failed to save patient"
// @Router /api/v1/dental/patient [post]
func CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient models.Patient
	if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if patient.ID == "" {
		patient.ID = uuid.NewString()
	}

	if err := patient.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if patient.CreatedAt == "" {
		patient.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if patient.UpdatedAt == "" {
		patient.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Patients"),
		Item: map[string]types.AttributeValue{
			"ID":           &types.AttributeValueMemberS{Value: patient.ID},
			"Name":         &types.AttributeValueMemberS{Value: patient.Name},
			"Email":        &types.AttributeValueMemberS{Value: patient.Email},
			"Phone":        &types.AttributeValueMemberS{Value: patient.Phone},
			"DateOfBirth":  &types.AttributeValueMemberS{Value: patient.DateOfBirth},
			"MedicalNotes": &types.AttributeValueMemberS{Value: patient.MedicalNotes},
			"CreatedAt":    &types.AttributeValueMemberS{Value: patient.CreatedAt},
			"UpdatedAt":    &types.AttributeValueMemberS{Value: patient.UpdatedAt},
		},
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})

	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Patient with this ID already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to save patient", http.StatusInternalServerError)
		log.Printf("Error saving patient: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(patient)
}

// GetAllPatients godoc
// @Summary Get all patients
// @Description Get a list of all patients
// @Tags patients
// @Produce json
// @Success 200 {array} models.Patient
// @Failure 500 {string} string "Failed to retrieve patients"
// @Router /api/v1/dental/patient [get]
func GetAllPatients(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Patients"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve patients", http.StatusInternalServerError)
		log.Printf("Error scanning patients: %v", err)
		return
	}

	var patients []models.Patient
	for _, item := range result.Items {
		var patient models.Patient
		if err := attributevalue.UnmarshalMap(item, &patient); err != nil {
			log.Printf("Error unmarshaling patient: %v", err)
			continue
		}
		patients = append(patients, patient)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// GetPatientByID godoc
// @Summary Get patient by ID
// @Description Get a patient by their ID
// @Tags patients
// @Produce json
// @Param id path string true "Patient ID"
// @Success 200 {object} models.Patient
// @Failure 404 {string} string "Patient not found"
// @Failure 500 {string} string "Failed to retrieve patient"
// @Router /api/v1/dental/patient/{id} [get]
func GetPatientByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Patients"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve patient", http.StatusInternalServerError)
		log.Printf("Error fetching patient with ID %s: %v", id, err)
		return
	}
	if result.Item == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	var patient models.Patient
	if err = attributevalue.UnmarshalMap(result.Item, &patient); err != nil {
		http.Error(w, "Failed to unmarshal patient data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling patient data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// GetPatientByName godoc
// @Summary Get patient by name
// @Description Get patients by their name (partial match)
// @Tags patients
// @Produce json
// @Param name path string true "Patient Name"
// @Success 200 {array} models.Patient
// @Failure 500 {string} string "Failed to retrieve patients"
// @Router /api/v1/dental/patient/name/{name} [get]
func GetPatientByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("Patients"),
		FilterExpression: aws.String("contains(#name, :name)"),
		ExpressionAttributeNames: map[string]string{
			"#name": "Name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve patients", http.StatusInternalServerError)
		log.Printf("Error scanning patients by name: %v", err)
		return
	}

	var patients []models.Patient
	for _, item := range result.Items {
		var patient models.Patient
		if err := attributevalue.UnmarshalMap(item, &patient); err != nil {
			log.Printf("Error unmarshaling patient: %v", err)
			continue
		}
		patients = append(patients, patient)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// UpdatePatient godoc
// @Summary Update an existing patient
// @Description Update fields of an existing patient by providing their ID
// @Tags patients
// @Accept json
// @Produce json
// @Param id path string true "Patient ID"
// @Param patient body models.Patient true "Patient data (ID will be ignored)"
// @Success 200 {object} models.Patient
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 404 {string} string "Patient not found"
// @Failure 500 {string} string "Failed to update patient"
// @Router /api/v1/dental/patient/{id} [put]
func UpdatePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Patients"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve patient", http.StatusInternalServerError)
		log.Printf("Error fetching patient with ID %s: %v", id, err)
		return
	}
	if result.Item == nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	var currentPatient models.Patient
	if err = attributevalue.UnmarshalMap(result.Item, &currentPatient); err != nil {
		http.Error(w, "Failed to unmarshal patient data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling patient data: %v", err)
		return
	}

	var updatedData models.Patient
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedData.Name != "" {
		currentPatient.Name = updatedData.Name
	}
	if updatedData.Email != "" {
		currentPatient.Email = updatedData.Email
	}
	if updatedData.Phone != "" {
		currentPatient.Phone = updatedData.Phone
	}
	if updatedData.DateOfBirth != "" {
		currentPatient.DateOfBirth = updatedData.DateOfBirth
	}
	if updatedData.MedicalNotes != "" {
		currentPatient.MedicalNotes = updatedData.MedicalNotes
	}

	if err := currentPatient.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentPatient.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	_, err = config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Patients"),
		Item: map[string]types.AttributeValue{
			"ID":           &types.AttributeValueMemberS{Value: currentPatient.ID},
			"Name":         &types.AttributeValueMemberS{Value: currentPatient.Name},
			"Email":        &types.AttributeValueMemberS{Value: currentPatient.Email},
			"Phone":        &types.AttributeValueMemberS{Value: currentPatient.Phone},
			"DateOfBirth":  &types.AttributeValueMemberS{Value: currentPatient.DateOfBirth},
			"MedicalNotes": &types.AttributeValueMemberS{Value: currentPatient.MedicalNotes},
			"CreatedAt":    &types.AttributeValueMemberS{Value: currentPatient.CreatedAt},
			"UpdatedAt":    &types.AttributeValueMemberS{Value: currentPatient.UpdatedAt},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update patient", http.StatusInternalServerError)
		log.Printf("Error updating patient: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentPatient)
}

// DeletePatient godoc
// @Summary Delete a patient
// @Description Delete a patient by their ID
// @Tags patients
// @Param id path string true "Patient ID"
// @Success 204 "Patient deleted successfully"
// @Failure 404 {string} string "Patient not found"
// @Failure 500 {string} string "Failed to delete patient"
// @Router /api/v1/dental/patient/{id} [delete]
func DeletePatient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := config.DBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Patients"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete patient", http.StatusInternalServerError)
		log.Printf("Error deleting patient: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}