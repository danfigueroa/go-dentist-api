package main

import (
	"log"
	"net/http"

	_ "go-dentist-api/docs"
	"go-dentist-api/internal/adapters/primary/http/handlers"
	"go-dentist-api/internal/adapters/primary/http/routes"
	"go-dentist-api/internal/adapters/secondary/database"
	"go-dentist-api/internal/application/usecases"
	"go-dentist-api/shared/config"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Dental SaaS API
// @version 1.0
// @description Multi-module SaaS platform for dental clinic management, including dental operations and financial management.

// @contact.name Daniel Figueroa API Support
// @contact.email danielmfigueroa@gmail.com

// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize DynamoDB
	config.InitDynamoDB()

	// Initialize repositories with DynamoDB client
	clinicRepo := database.NewDynamoDBClinicRepository(config.DBClient, "Clinics")
	dentistRepo := database.NewDynamoDBDentistRepository(config.DBClient, "Dentists")

	dentistUseCase := usecases.NewDentistUseCase(dentistRepo, clinicRepo)

	dentistHandler := handlers.NewDentistHandler(dentistUseCase)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	routes.RegisterDentistRoutes(apiRouter, dentistHandler)

	log.Println("Dental SaaS running on http://localhost:8080")
	log.Println("API documentation available at http://localhost:8080/swagger/")
	log.Fatal(http.ListenAndServe(":8080", router))
}
