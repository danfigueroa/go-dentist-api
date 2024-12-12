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
// @Router /dentist [post]
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

type dynamoDentist struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CRO       string `json:"cro"`
	Country   string `json:"country"`
	Specialty string `json:"specialty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func dynamoToModel(dd dynamoDentist) (models.Dentist, error) {
	var dentist models.Dentist
	dentist.ID = dd.ID
	dentist.Name = dd.Name
	dentist.Email = dd.Email
	dentist.Phone = dd.Phone
	dentist.CRO = dd.CRO
	dentist.Country = dd.Country
	dentist.Specialty = dd.Specialty

	var err error
	if dd.CreatedAt != "" {
		dentist.CreatedAt, err = time.Parse(time.RFC3339, dd.CreatedAt)
		if err != nil {
			return dentist, err
		}
	}
	if dd.UpdatedAt != "" {
		dentist.UpdatedAt, err = time.Parse(time.RFC3339, dd.UpdatedAt)
		if err != nil {
			return dentist, err
		}
	}

	return dentist, nil
}

// GetAllDentists godoc
// @Summary Get all dentists
// @Description Retrieve all registered dentists
// @Tags dentists
// @Produce json
// @Success 200 {array} models.Dentist
// @Failure 500 {string} string "Failed to retrieve dentists"
// @Router /dentists [get]
func GetAllDentists(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Dentists"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentists", http.StatusInternalServerError)
		log.Printf("Error fetching dentists: %v", err)
		return
	}

	var dynamoDentists []dynamoDentist
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dynamoDentists)
	if err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	var dentists []models.Dentist
	for _, dd := range dynamoDentists {
		d, convErr := dynamoToModel(dd)
		if convErr != nil {
			http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
			log.Printf("Error parsing date fields: %v", convErr)
			return
		}
		dentists = append(dentists, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}

// GetDentistByID godoc
// @Summary Get dentist by ID
// @Description Retrieve a single dentist by providing its ID
// @Tags dentists
// @Produce json
// @Param id path string true "Dentist ID"
// @Success 200 {object} models.Dentist
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to retrieve dentist"
// @Router /dentist/{id} [get]
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

	var dd dynamoDentist
	err = attributevalue.UnmarshalMap(result.Item, &dd)
	if err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	dentist, convErr := dynamoToModel(dd)
	if convErr != nil {
		http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
		log.Printf("Error parsing date fields: %v", convErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentist)
}

// UpdateDentist godoc
// @Summary Update an existing dentist
// @Description Update fields of an existing dentist by providing its ID
// @Tags dentists
// @Accept json
// @Produce json
// @Param id path string true "Dentist ID"
// @Param dentist body models.Dentist true "Dentist data (ID will be ignored)"
// @Success 200 {object} models.Dentist
// @Failure 400 {string} string "Invalid request body or missing required fields"
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to update dentist"
// @Router /dentist/{id} [put]
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

	var dd dynamoDentist
	if err = attributevalue.UnmarshalMap(result.Item, &dd); err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	currentDentist, convErr := dynamoToModel(dd)
	if convErr != nil {
		http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
		log.Printf("Error parsing date fields: %v", convErr)
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

	// Valida campos obrigatórios após atualização
	if err := currentDentist.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	currentDentist.UpdatedAt = time.Now().UTC()

	createdAtStr := currentDentist.CreatedAt.Format(time.RFC3339)
	updatedAtStr := currentDentist.UpdatedAt.Format(time.RFC3339)

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
// @Summary Delete a dentist by ID
// @Description Delete a single dentist by providing its ID
// @Tags dentists
// @Produce json
// @Param id path string true "Dentist ID"
// @Success 204 "No Content"
// @Failure 404 {string} string "Dentist not found"
// @Failure 500 {string} string "Failed to delete dentist"
// @Router /dentist/{id} [delete]
func DeleteDentist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := config.DBClient.DeleteItem(r.Context(), &dynamodb.DeleteItemInput{
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

// GetDentistByName godoc
// @Summary Get dentist by name
// @Description Retrieve dentists by providing a name (partial match)
// @Tags dentists
// @Produce json
// @Param name path string true "Dentist Name"
// @Success 200 {array} models.Dentist
// @Failure 404 {string} string "No dentists found with this name"
// @Failure 500 {string} string "Failed to retrieve dentists"
// @Router /dentist/name/{name} [get]
func GetDentistByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Buscar todos os dentistas
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Dentists"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentists", http.StatusInternalServerError)
		log.Printf("Error fetching dentists: %v", err)
		return
	}

	var dynamoDentists []dynamoDentist
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dynamoDentists)
	if err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	// Filtrar dentistas pelo nome (correspondência parcial)
	var matchingDentists []models.Dentist
	for _, dd := range dynamoDentists {
		// Verificar se o nome contém a string de busca (case insensitive)
		if strings.Contains(strings.ToLower(dd.Name), strings.ToLower(name)) {
			d, convErr := dynamoToModel(dd)
			if convErr != nil {
				http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
				log.Printf("Error parsing date fields: %v", convErr)
				return
			}
			matchingDentists = append(matchingDentists, d)
		}
	}

	if len(matchingDentists) == 0 {
		http.Error(w, "No dentists found with this name", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matchingDentists)
}

// GetDentistByCRO godoc
// @Summary Get dentist by CRO
// @Description Retrieve a dentist by providing its CRO
// @Tags dentists
// @Produce json
// @Param cro path string true "Dentist CRO"
// @Success 200 {object} models.Dentist
// @Failure 404 {string} string "No dentist found with this CRO"
// @Failure 500 {string} string "Failed to retrieve dentist"
// @Router /dentist/cro/{cro} [get]
func GetDentistByCRO(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cro := vars["cro"]

	// Buscar todos os dentistas
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Dentists"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentists", http.StatusInternalServerError)
		log.Printf("Error fetching dentists: %v", err)
		return
	}

	var dynamoDentists []dynamoDentist
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dynamoDentists)
	if err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	// Encontrar o dentista com o CRO correspondente
	for _, dd := range dynamoDentists {
		if strings.EqualFold(dd.CRO, cro) {
			dentist, convErr := dynamoToModel(dd)
			if convErr != nil {
				http.Error(w, "Failed to parse date fields", http.StatusInternalServerError)
				log.Printf("Error parsing date fields: %v", convErr)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(dentist)
			return
		}
	}

	http.Error(w, "No dentist found with this CRO", http.StatusNotFound)
}
