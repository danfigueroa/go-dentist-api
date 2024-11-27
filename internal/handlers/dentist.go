package handlers

import (
	"context"
	"encoding/json"
	"go-dentist-api/internal/config"
	"go-dentist-api/internal/models"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gorilla/mux"
)

func CreateDentist(w http.ResponseWriter, r *http.Request) {
	var dentist models.Dentist
	if err := json.NewDecoder(r.Body).Decode(&dentist); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := dentist.IsValid(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
			"CreatedAt": &types.AttributeValueMemberS{Value: dentist.CreatedAt},
			"UpdatedAt": &types.AttributeValueMemberS{Value: dentist.UpdatedAt},
		},
	})
	if err != nil {
		http.Error(w, "Failed to save dentist", http.StatusInternalServerError)
		log.Printf("Error saving dentist: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dentist)
}

func GetAllDentists(w http.ResponseWriter, r *http.Request) {
	result, err := config.DBClient.Scan(r.Context(), &dynamodb.ScanInput{
		TableName: aws.String("Dentists"),
	})
	if err != nil {
		http.Error(w, "Failed to retrieve dentists", http.StatusInternalServerError)
		log.Printf("Error fetching dentists: %v", err)
		return
	}

	var dentists []models.Dentist
	for _, item := range result.Items {
		var dentist models.Dentist
		err := attributevalue.UnmarshalMap(item, &dentist)
		if err != nil {
			http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
			log.Printf("Error unmarshaling dentist data: %v", err)
			return
		}
		dentists = append(dentists, dentist)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentists)
}

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
	err = attributevalue.UnmarshalMap(result.Item, &dentist)
	if err != nil {
		http.Error(w, "Failed to unmarshal dentist data", http.StatusInternalServerError)
		log.Printf("Error unmarshaling dentist data: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dentist)
}
