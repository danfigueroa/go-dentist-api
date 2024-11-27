package main

import (
	"go-dentist-api/internal/config"
	"go-dentist-api/internal/router"
	"log"
	"net/http"
)

func main() {
	config.InitDynamoDB()

	r := router.SetupRouter()

	port := "8080"
	log.Printf("Server running on http://localhost:%s", port)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
