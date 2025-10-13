package routes

import (
	"go-dentist-api/internal/adapters/primary/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterDentistRoutes(router *mux.Router, dentistHandler *handlers.DentistHandler) {
	dentistRouter := router.PathPrefix("/dentists").Subrouter()

	dentistRouter.HandleFunc("", dentistHandler.CreateDentist).Methods("POST")
	dentistRouter.HandleFunc("", dentistHandler.ListDentists).Methods("GET")
	dentistRouter.HandleFunc("/search", dentistHandler.GetDentistByCRO).Methods("GET").Queries("cro", "{cro}")
	dentistRouter.HandleFunc("/search", dentistHandler.GetDentistsByName).Methods("GET").Queries("name", "{name}")
	dentistRouter.HandleFunc("/{id}", dentistHandler.GetDentist).Methods("GET")
	dentistRouter.HandleFunc("/{id}", dentistHandler.UpdateDentist).Methods("PUT")
	dentistRouter.HandleFunc("/{id}", dentistHandler.DeleteDentist).Methods("DELETE")
	dentistRouter.HandleFunc("/clinic/{clinicId}", dentistHandler.GetDentistsByClinic).Methods("GET")
}
