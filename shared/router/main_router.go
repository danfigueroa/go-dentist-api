package router

import (
	"dental-saas/modules/dental/router"
	"net/http"

	"github.com/gorilla/mux"
)

// NewMainRouter creates the main router that orchestrates all module routers
func NewMainRouter() *mux.Router {
	mainRouter := mux.NewRouter()

	// Health check endpoint
	mainRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"dental-saas"}`))
	}).Methods("GET")

	// API version info
	mainRouter.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version":"1.0","modules":["dental","financial"]}`))
	}).Methods("GET")

	// Register dental module routes
	dentalRouter := router.NewDentalRouter()
	mainRouter.PathPrefix("/api/v1/dental").Handler(dentalRouter)

	// TODO: Register financial module routes when implemented
	// financialRouter := financial_router.NewFinancialRouter()
	// mainRouter.PathPrefix("/api/v1/financial").Handler(financialRouter)

	// TODO: Register other future modules here

	return mainRouter
}