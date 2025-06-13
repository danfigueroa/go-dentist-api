package main

import (
	"log"
	"net/http"

	_ "dental-saas/docs"
	"dental-saas/shared/config"
	"dental-saas/shared/router"

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
	config.InitDynamoDB()

	r := router.NewMainRouter()

	// Adiciona o Swagger na rota principal
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Dental SaaS running on http://localhost:8080")
	log.Println("API documentation available at http://localhost:8080/swagger/")
	log.Fatal(http.ListenAndServe(":8080", r))
}
