package main

import (
	"log"
	"net/http"

	_ "go-dentist-api/docs"
	"go-dentist-api/internal/config"
	"go-dentist-api/internal/router"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Dentist API
// @version 1.0
// @description This is a sample server for managing dentists.

// @Daniel Figueroa API Support
// @contact.email danielmfigueroa@gmail.com

// @host localhost:8080
// @BasePath /
func main() {
	config.InitDynamoDB()

	r := router.SetupRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
