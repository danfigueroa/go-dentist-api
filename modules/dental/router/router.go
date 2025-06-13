package router

import (
	"dental-saas/modules/dental/handlers"

	"github.com/gorilla/mux"
)

// NewDentalRouter creates and configures routes for the dental module
func NewDentalRouter() *mux.Router {
	r := mux.NewRouter()

	// Create a subrouter for dental module with /api/v1/dental prefix
	dentalRouter := r.PathPrefix("/api/v1/dental").Subrouter()

	// Dentist routes
	dentalRouter.HandleFunc("/dentist", handlers.CreateDentist).Methods("POST")
	dentalRouter.HandleFunc("/dentist", handlers.GetAllDentists).Methods("GET")
	dentalRouter.HandleFunc("/dentist/name/{name}", handlers.GetDentistByName).Methods("GET")
	dentalRouter.HandleFunc("/dentist/cro/{cro}", handlers.GetDentistByCRO).Methods("GET")
	dentalRouter.HandleFunc("/dentist/{id}", handlers.GetDentistByID).Methods("GET")
	dentalRouter.HandleFunc("/dentist/{id}", handlers.UpdateDentist).Methods("PUT")
	dentalRouter.HandleFunc("/dentist/{id}", handlers.DeleteDentist).Methods("DELETE")

	// Patient routes
	dentalRouter.HandleFunc("/patient", handlers.CreatePatient).Methods("POST")
	dentalRouter.HandleFunc("/patient", handlers.GetAllPatients).Methods("GET")
	dentalRouter.HandleFunc("/patient/{id}", handlers.GetPatientByID).Methods("GET")
	dentalRouter.HandleFunc("/patient/name/{name}", handlers.GetPatientByName).Methods("GET")
	dentalRouter.HandleFunc("/patient/{id}", handlers.UpdatePatient).Methods("PUT")
	dentalRouter.HandleFunc("/patient/{id}", handlers.DeletePatient).Methods("DELETE")

	// Procedure routes
	dentalRouter.HandleFunc("/procedure", handlers.CreateProcedure).Methods("POST")
	dentalRouter.HandleFunc("/procedure", handlers.GetAllProcedures).Methods("GET")
	dentalRouter.HandleFunc("/procedure/{id}", handlers.GetProcedureByID).Methods("GET")
	dentalRouter.HandleFunc("/procedure/name/{name}", handlers.GetProcedureByName).Methods("GET")
	dentalRouter.HandleFunc("/procedure/{id}", handlers.UpdateProcedure).Methods("PUT")
	dentalRouter.HandleFunc("/procedure/{id}", handlers.DeleteProcedure).Methods("DELETE")

	// Appointment routes
	dentalRouter.HandleFunc("/appointment", handlers.CreateAppointment).Methods("POST")
	dentalRouter.HandleFunc("/appointment", handlers.GetAllAppointments).Methods("GET")
	dentalRouter.HandleFunc("/appointment/{id}", handlers.GetAppointmentByID).Methods("GET")
	dentalRouter.HandleFunc("/appointment/patient/{patientId}", handlers.GetAppointmentsByPatient).Methods("GET")
	dentalRouter.HandleFunc("/appointment/dentist/{dentistId}", handlers.GetAppointmentsByDentist).Methods("GET")
	dentalRouter.HandleFunc("/appointment/{id}", handlers.UpdateAppointment).Methods("PUT")
	dentalRouter.HandleFunc("/appointment/{id}", handlers.DeleteAppointment).Methods("DELETE")

	return r
}