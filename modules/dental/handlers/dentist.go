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

// CreateDentist godoc
// @Summary Create a new dentist
// @Description Create a new dentist by providing the details
// @Tags dentists
// @Accept json
// @Produce json
// @Param dentist body models.Dentist true "Dentist data"
// @Success 201 {object} models.Dentist
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 409 {string} string "Dentist with this ID already exists"
// @Failure 500 {string} string "Failed to save dentist"
// @Router /api/v1/dental/dentist [post]
func CreateDentist(w http.ResponseWriter, r *http.Request) {
	var dentist models.Dentist
	if err := json.NewDecoder(r.Body).Decode(&dentist); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if dentist.ID == "" {
		dentist.ID = uuid.NewString()
	}

	if err := dentist.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if dentist.CreatedAt.IsZero() {
		dentist.CreatedAt = time.Now().UTC()
	}
	if dentist.UpdatedAt.IsZero() {
		dentist.UpdatedAt = time.Now().UTC()
	}

	createdAtStr := dentist.CreatedAt.Format(time.RFC3339)
	updatedAtStr := dentist.UpdatedAt.Format(time.RFC3339)

	_, err := config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Dentists"),
		Item: map[string]types.AttributeValue{
			"ID":        &types.AttributeValueMemberS{Value: dentist.ID},
			"Name":      &types.AttributeValueMemberS{Value: dentist.Name},
			"Email":     &types.AttributeValueMemberS{Value: dentist.Email},
			"Phone":     &types.AttributeValueMemberS{Value: dentist.Phone},
			"CRO":       &types.AttributeValueMemberS{Value: dentist.CRO},
			"Country":   &types.AttributeValueMemberS{Value: dentist.Country},
			"Specialty": &types.AttributeValueMemberS{Value: dentist.Specialty},
			"CreatedAt": &types.AttributeValueMemberS{Value: createdAtStr},
			"UpdatedAt": &types.AttributeValueMemberS{Value: updatedAtStr},
		},
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})

	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Dentist with this ID already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to save dentist", http.StatusInternalServerError)
		log.Printf("Error saving dentist: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dentist)
}

// GetAllDentists godoc
// @Summary Get all dentists
// @Description Get a list of all dentists
// @Tags dentists
// @Produce json
// @Success 200 {array} models.Dentist
// @Failure 500 {string} string "Failed to retrieve dentists"
// @Router /api/v1/dental/dentist [get]
func GetAllDentists(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName: aws.String("Dentists"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentists", http.StatusInternalServerError)
		log.Printf("Error scanning dentists: %v", err)
		return
	}

	var dentists []models.Dentist
	for _, item := range result.Items {
		var dentist models.Dentist
		if err := attributevalue.UnmarshalMap(item, &dentist); err != nil {
			log.Printf("Error unmarshaling dentist: %v", err)
			continue
		}
		dentists = append(dentists, dentist)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}

// GetDentistByID godoc
// @Summary Get dentist by ID
// @Description Get a dentist by their ID
// @Tags dentists
// @Produce json
// @Param id path string true "Dentist ID"
// @Success 200 {object} models.Dentist
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to retrieve dentist"
// @Router /api/v1/dental/dentist/{id} [get]
func GetDentistByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Dentists"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentist", http.StatusInternalServerError)
		log.Printf("Error fetching dentist with ID %s: %v", id, err)
		return
	}
	if result.Item == nil {
		http.Error(w, "Dentist not found", http.StatusNotFound)
		return
	}

	var dentist models.Dentist
	if err = attributevalue.UnmarshalMap(result.Item, &dentist); err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentist)
}

// GetDentistByName godoc
// @Summary Get dentist by name
// @Description Get dentists by their name (partial match)
// @Tags dentists
// @Produce json
// @Param name path string true "Dentist Name"
// @Success 200 {array} models.Dentist
// @Failure 500 {string} string "Failed to retrieve dentists"
// @Router /api/v1/dental/dentist/name/{name} [get]
func GetDentistByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("Dentists"),
		FilterExpression: aws.String("contains(#name, :name)"),
		ExpressionAttributeNames: map[string]string{
			"#name": "Name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentists", http.StatusInternalServerError)
		log.Printf("Error scanning dentists by name: %v", err)
		return
	}

	var dentists []models.Dentist
	for _, item := range result.Items {
		var dentist models.Dentist
		if err := attributevalue.UnmarshalMap(item, &dentist); err != nil {
			log.Printf("Error unmarshaling dentist: %v", err)
			continue
		}
		dentists = append(dentists, dentist)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}

// GetDentistByCRO godoc
// @Summary Get dentist by CRO
// @Description Get a dentist by their CRO number
// @Tags dentists
// @Produce json
// @Param cro path string true "Dentist CRO"
// @Success 200 {object} models.Dentist
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to retrieve dentist"
// @Router /api/v1/dental/dentist/cro/{cro} [get]
func GetDentistByCRO(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cro := vars["cro"]

	result, err := config.DBClient.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:        aws.String("Dentists"),
		FilterExpression: aws.String("CRO = :cro"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cro": &types.AttributeValueMemberS{Value: cro},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentist", http.StatusInternalServerError)
		log.Printf("Error scanning dentist by CRO: %v", err)
		return
	}

	if len(result.Items) == 0 {
		http.Error(w, "Dentist not found", http.StatusNotFound)
		return
	}

	var dentist models.Dentist
	if err = attributevalue.UnmarshalMap(result.Items[0], &dentist); err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentist)
}

// UpdateDentist godoc
// @Summary Update dentist
// @Description Update an existing dentist
// @Tags dentists
// @Accept json
// @Produce json
// @Param id path string true "Dentist ID"
// @Param dentist body models.Dentist true "Updated dentist data"
// @Success 200 {object} models.Dentist
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to update dentist"
// @Router /api/v1/dental/dentist/{id} [put]
func UpdateDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := config.DBClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("Dentists"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentist", http.StatusInternalServerError)
		log.Printf("Error fetching dentist with ID %s: %v", id, err)
		return
	}
	if result.Item == nil {
		http.Error(w, "Dentist not found", http.StatusNotFound)
		return
	}

	var currentDentist models.Dentist
	if err = attributevalue.UnmarshalMap(result.Item, &currentDentist); err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	var updatedData models.Dentist
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if updatedData.Name != "" {
		currentDentist.Name = updatedData.Name
	}
	if updatedData.Email != "" {
		currentDentist.Email = updatedData.Email
	}
	if updatedData.Phone != "" {
		currentDentist.Phone = updatedData.Phone
	}
	if updatedData.CRO != "" {
		currentDentist.CRO = updatedData.CRO
	}
	if updatedData.Country != "" {
		currentDentist.Country = updatedData.Country
	}
	if updatedData.Specialty != "" {
		currentDentist.Specialty = updatedData.Specialty
	}

	if err := currentDentist.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDentist.UpdatedAt = time.Now().UTC()
	updatedAtStr := currentDentist.UpdatedAt.Format(time.RFC3339)
	createdAtStr := currentDentist.CreatedAt.Format(time.RFC3339)

	_, err = config.DBClient.PutItem(r.Context(), &dynamodb.PutItemInput{
		TableName: aws.String("Dentists"),
		Item: map[string]types.AttributeValue{
			"ID":        &types.AttributeValueMemberS{Value: currentDentist.ID},
			"Name":      &types.AttributeValueMemberS{Value: currentDentist.Name},
			"Email":     &types.AttributeValueMemberS{Value: currentDentist.Email},
			"Phone":     &types.AttributeValueMemberS{Value: currentDentist.Phone},
			"CRO":       &types.AttributeValueMemberS{Value: currentDentist.CRO},
			"Country":   &types.AttributeValueMemberS{Value: currentDentist.Country},
			"Specialty": &types.AttributeValueMemberS{Value: currentDentist.Specialty},
			"CreatedAt": &types.AttributeValueMemberS{Value: createdAtStr},
			"UpdatedAt": &types.AttributeValueMemberS{Value: updatedAtStr},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Dentist not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to update dentist", http.StatusInternalServerError)
		log.Printf("Error updating dentist: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentDentist)
}

// DeleteDentist godoc
// @Summary Delete dentist
// @Description Delete a dentist by ID
// @Tags dentists
// @Param id path string true "Dentist ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to delete dentist"
// @Router /api/v1/dental/dentist/{id} [delete]
func DeleteDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := config.DBClient.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("Dentists"),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})
	if err != nil {
		var cfe *types.ConditionalCheckFailedException
		if errors.As(err, &cfe) {
			http.Error(w, "Dentist not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete dentist", http.StatusInternalServerError)
		log.Printf("Error deleting dentist: %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}