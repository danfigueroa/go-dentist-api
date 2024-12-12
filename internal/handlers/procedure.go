package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-dentist-api/internal/config"
	"go-dentist-api/internal/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type dynamoProcedure struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	DentistID    string `json:"dentist_id"`
	PatientID    string `json:"patient_id"`
	PerformedAt  string `json:"performed_at"`
	Observations string `json:"observations"`
	Cost         string `json:"cost"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateProcedure godoc
// @Summary Create a new procedure
// @Description Create a new procedure by providing the details
// @Tags procedures
// @Accept json
// @Produce json
// @Param procedure body models.Procedure true "Procedure data"
// @Success 201 {object} models.Procedure
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 409 {string} string "Procedure with this ID already exists"
// @Failure 500 {string} string "Failed to save procedure"
// @Router /procedure [post]
func CreateProcedure(w http.ResponseWriter, r *http.Request) {
	var proc models.Procedure
	if err := json.NewDecoder(r.Body).Decode(&proc); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if proc.ID == "" {
		proc.ID = uuid.NewString()
	}

	if err := proc.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if proc.CreatedAt.IsZero() {
		proc.CreatedAt = time.Now().UTC()
	}
	if proc.UpdatedAt.IsZero() {
		proc.UpdatedAt = time.Now().UTC()
	}

	createdAtStr := proc.CreatedAt.Format(time.RFC3339)
	updatedAtStr := proc.UpdatedAt.Format(time.RFC3339)
	performedAtStr := proc.PerformedAt.Format(time.RFC3339)
	costStr := fmt.Sprintf("%f", proc.Cost)

	_, err := config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Procedures"),
		Item: map[string]types.AttributeValue{
			"ID":           &types.AttributeValueMemberS{Value: proc.ID},
			"Name":         &types.AttributeValueMemberS{Value: proc.Name},
			"Type":         &types.AttributeValueMemberS{Value: proc.Type},
			"DentistID":    &types.AttributeValueMemberS{Value: proc.DentistID},
			"PatientID":    &types.AttributeValueMemberS{Value: proc.PatientID},
			"PerformedAt":  &types.AttributeValueMemberS{Value: performedAtStr},
			"Observations": &types.AttributeValueMemberS{Value: proc.Observations},
			"Cost":         &types.AttributeValueMemberN{Value: costStr}, // Cost como número
			"CreatedAt":    &types.AttributeValueMemberS{Value: createdAtStr},
			"UpdatedAt":    &types.AttributeValueMemberS{Value: updatedAtStr},
		},
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Procedure with this ID already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to save procedure", http.StatusInternalServerError)
		log.Printf("Error saving procedure: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(proc)
}

// GetAllProcedures godoc
// @Summary Get all procedures
// @Description Retrieve all registered procedures
// @Tags procedures
// @Produce json
// @Success 200 {array} models.Procedure
// @Failure 500 {string} string "Failed to retrieve procedures"
// @Router /procedures [get]
func GetAllProcedures(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Procedures"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve procedures", http.StatusInternalServerError)
		log.Printf("Error fetching procedures: %v", err)
		return
	}

	var dynamoProcs []dynamoProcedure
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dynamoProcs)
	if err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	var procs []models.Procedure
	for _, dp := range dynamoProcs {
		p, convErr := dynamoToProcedure(dp)
		if convErr != nil {
			http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
			log.Printf("Error parsing date fields: %v", convErr)
			return
		}
		procs = append(procs, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procs)
}

// GetProcedureByID godoc
// @Summary Get procedure by ID
// @Description Retrieve a single procedure by providing its ID
// @Tags procedures
// @Produce json
// @Param id path string true "Procedure ID"
// @Success 200 {object} models.Procedure
// @Failure 404 {string} string "Procedure not found"
// @Failure 500 {string} string "Failed to retrieve procedure"
// @Router /procedure/{id} [get]
func GetProcedureByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Procedures"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve procedure", http.StatusInternalServerError)
		log.Printf("Error fetching procedure with ID %s: %v", id, err)
		return
	}

	if result.Item == nil {
		http.Error(w, "Procedure not found", http.StatusNotFound)
		return
	}

	var dp dynamoProcedure
	err = attributevalue.UnmarshalMap(result.Item, &dp)
	if err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	proc, convErr := dynamoToProcedure(dp)
	if convErr != nil {
		http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
		log.Printf("Error parsing date fields: %v", convErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proc)
}

// UpdateProcedure godoc
// @Summary Update an existing procedure
// @Description Update fields of an existing procedure by providing its ID
// @Tags procedures
// @Accept json
// @Produce json
// @Param id path string true "Procedure ID"
// @Param procedure body models.Procedure true "Procedure data (ID ignored)"
// @Success 200 {object} models.Procedure
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 404 {string} string "Procedure not found"
// @Failure 500 {string} string "Failed to update procedure"
// @Router /procedure/{id} [put]
func UpdateProcedure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Procedures"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve procedure", http.StatusInternalServerError)
		log.Printf("Error fetching procedure with ID %s: %v", id, err)
		return
	}
	if result.Item == nil {
		http.Error(w, "Procedure not found", http.StatusNotFound)
		return
	}

	var dp dynamoProcedure
	if err = attributevalue.UnmarshalMap(result.Item, &dp); err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	currentProc, convErr := dynamoToProcedure(dp)
	if convErr != nil {
		http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
		log.Printf("Error parsing date fields: %v", convErr)
		return
	}

	var updatedData models.Procedure
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedData.Name != "" {
		currentProc.Name = updatedData.Name
	}
	if updatedData.Type != "" {
		currentProc.Type = updatedData.Type
	}
	if updatedData.DentistID != "" {
		currentProc.DentistID = updatedData.DentistID
	}
	if !updatedData.PerformedAt.IsZero() {
		currentProc.PerformedAt = updatedData.PerformedAt
	}
	if updatedData.Observations != "" {
		currentProc.Observations = updatedData.Observations
	}
	if updatedData.PatientID != "" {
		currentProc.PatientID = updatedData.PatientID
	}
	if updatedData.Cost != 0 {
		currentProc.Cost = updatedData.Cost
	}

	if err := currentProc.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentProc.UpdatedAt = time.Now().UTC()

	createdAtStr := currentProc.CreatedAt.Format(time.RFC3339)
	updatedAtStr := currentProc.UpdatedAt.Format(time.RFC3339)
	performedAtStr := currentProc.PerformedAt.Format(time.RFC3339)
	costStr := fmt.Sprintf("%f", currentProc.Cost)

	_, err = config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Procedures"),
		Item: map[string]types.AttributeValue{
			"ID":           &types.AttributeValueMemberS{Value: currentProc.ID},
			"Name":         &types.AttributeValueMemberS{Value: currentProc.Name},
			"Type":         &types.AttributeValueMemberS{Value: currentProc.Type},
			"DentistID":    &types.AttributeValueMemberS{Value: currentProc.DentistID},
			"PatientID":    &types.AttributeValueMemberS{Value: currentProc.PatientID},
			"PerformedAt":  &types.AttributeValueMemberS{Value: performedAtStr},
			"Observations": &types.AttributeValueMemberS{Value: currentProc.Observations},
			"Cost":         &types.AttributeValueMemberN{Value: costStr},
			"CreatedAt":    &types.AttributeValueMemberS{Value: createdAtStr},
			"UpdatedAt":    &types.AttributeValueMemberS{Value: updatedAtStr},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Procedure not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update procedure", http.StatusInternalServerError)
		log.Printf("Error updating procedure: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentProc)
}

// DeleteProcedure godoc
// @Summary Delete a procedure by ID
// @Description Delete a single procedure by providing its ID
// @Tags procedures
// @Produce json
// @Param id path string true "Procedure ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Procedure not found"
// @Failure 500 {string} string "Failed to delete procedure"
// @Router /procedure/{id} [delete]
func DeleteProcedure(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := config.DBClient.DeleteItem(r.Context(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Procedures"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Procedure not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete procedure", http.StatusInternalServerError)
		log.Printf("Error deleting procedure: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func dynamoToProcedure(dp dynamoProcedure) (models.Procedure, error) {
	var proc models.Procedure
	proc.ID = dp.ID
	proc.Name = dp.Name
	proc.Type = dp.Type
	proc.DentistID = dp.DentistID
	proc.PatientID = dp.PatientID
	proc.Observations = dp.Observations

	var err error
	if dp.PerformedAt != "" {
		proc.PerformedAt, err = time.Parse(time.RFC3339, dp.PerformedAt)
		if err != nil {
			return proc, err
		}
	}
	if dp.CreatedAt != "" {
		proc.CreatedAt, err = time.Parse(time.RFC3339, dp.CreatedAt)
		if err != nil {
			return proc, err
		}
	}
	if dp.UpdatedAt != "" {
		proc.UpdatedAt, err = time.Parse(time.RFC3339, dp.UpdatedAt)
		if err != nil {
			return proc, err
		}
	}
	if dp.Cost != "" {
		costValue, cErr := strconv.ParseFloat(dp.Cost, 64)
		if cErr != nil {
			return proc, cErr
		}
		proc.Cost = costValue
	}

	return proc, nil
}
