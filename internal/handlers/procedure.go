package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"go-dentist-api/internal/config"
	"go-dentist-api/internal/models"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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
	var procedure models.Procedure
	if err := json.NewDecoder(r.Body).Decode(&procedure); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if procedure.ID == "" {
		procedure.ID = uuid.NewString()
	}

	if err := procedure.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if procedure.CreatedAt == "" {
		procedure.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if procedure.UpdatedAt == "" {
		procedure.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Procedures"),
		Item: map[string]types.AttributeValue{
			"ID":          &types.AttributeValueMemberS{Value: procedure.ID},
			"Name":        &types.AttributeValueMemberS{Value: procedure.Name},
			"Description": &types.AttributeValueMemberS{Value: procedure.Description},
			"Price":       &types.AttributeValueMemberS{Value: procedure.Price},
			"Duration":    &types.AttributeValueMemberS{Value: procedure.Duration},
			"CreatedAt":   &types.AttributeValueMemberS{Value: procedure.CreatedAt},
			"UpdatedAt":   &types.AttributeValueMemberS{Value: procedure.UpdatedAt},
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
	json.NewEncoder(w).Encode(procedure)
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

	var procedures []models.Procedure
	err = attributevalue.UnmarshalListOfMaps(result.Items, &procedures)
	if err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procedures)
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

	var procedure models.Procedure
	err = attributevalue.UnmarshalMap(result.Item, &procedure)
	if err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procedure)
}

// UpdateProcedure godoc
// @Summary Update an existing procedure
// @Description Update fields of an existing procedure by providing its ID
// @Tags procedures
// @Accept json
// @Produce json
// @Param id path string true "Procedure ID"
// @Param procedure body models.Procedure true "Procedure data (ID will be ignored)"
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

	var currentProcedure models.Procedure
	if err = attributevalue.UnmarshalMap(result.Item, &currentProcedure); err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	var updatedData models.Procedure
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedData.Name != "" {
		currentProcedure.Name = updatedData.Name
	}
	if updatedData.Description != "" {
		currentProcedure.Description = updatedData.Description
	}
	if updatedData.Price != "" {
		currentProcedure.Price = updatedData.Price
	}
	if updatedData.Duration != "" {
		currentProcedure.Duration = updatedData.Duration
	}

	// Valida campos obrigatórios após atualização
	if err := currentProcedure.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentProcedure.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	_, err = config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Procedures"),
		Item: map[string]types.AttributeValue{
			"ID":          &types.AttributeValueMemberS{Value: currentProcedure.ID},
			"Name":        &types.AttributeValueMemberS{Value: currentProcedure.Name},
			"Description": &types.AttributeValueMemberS{Value: currentProcedure.Description},
			"Price":       &types.AttributeValueMemberS{Value: currentProcedure.Price},
			"Duration":    &types.AttributeValueMemberS{Value: currentProcedure.Duration},
			"CreatedAt":   &types.AttributeValueMemberS{Value: currentProcedure.CreatedAt},
			"UpdatedAt":   &types.AttributeValueMemberS{Value: currentProcedure.UpdatedAt},
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
	json.NewEncoder(w).Encode(currentProcedure)
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

// GetProcedureByName godoc
// @Summary Get procedures by name
// @Description Retrieve procedures by providing a name (partial match)
// @Tags procedures
// @Produce json
// @Param name path string true "Procedure Name"
// @Success 200 {array} models.Procedure
// @Failure 404 {string} string "No procedures found with this name"
// @Failure 500 {string} string "Failed to retrieve procedures"
// @Router /procedure/name/{name} [get]
func GetProcedureByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Buscar todos os procedimentos
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Procedures"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve procedures", http.StatusInternalServerError)
		log.Printf("Error fetching procedures: %v", err)
		return
	}

	var procedures []map[string]types.AttributeValue
	for _, item := range result.Items {
		// Verificar se o nome contém a string de busca (case insensitive)
		if nameAttr, ok := item["Name"]; ok {
			if nameValue, ok := nameAttr.(*types.AttributeValueMemberS); ok {
				if strings.Contains(strings.ToLower(nameValue.Value), strings.ToLower(name)) {
					procedures = append(procedures, item)
				}
			}
		}
	}

	if len(procedures) == 0 {
		http.Error(w, "No procedures found with this name", http.StatusNotFound)
		return
	}

	var procedureList []models.Procedure
	err = attributevalue.UnmarshalListOfMaps(procedures, &procedureList)
	if err != nil {
		http.Error(w, "Failed to unmarshal procedure data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling procedure data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(procedureList)
}
