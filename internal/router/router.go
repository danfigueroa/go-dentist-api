package router

import (
	"go-dentist-api/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	// Dentist routes
	r.HandleFunc("/dentist", handlers.CreateDentist).Methods("POST")
	r.HandleFunc("/dentist", handlers.GetAllDentists).Methods("GET")
	r.HandleFunc("/dentist/name/{name}", handlers.GetDentistByName).Methods("GET")
	r.HandleFunc("/dentist/cro/{cro}", handlers.GetDentistByCRO).Methods("GET")
	r.HandleFunc("/dentist/{id}", handlers.GetDentistByID).Methods("GET")
	r.HandleFunc("/dentist/{id}", handlers.UpdateDentist).Methods("PUT")
	r.HandleFunc("/dentist/{id}", handlers.DeleteDentist).Methods("DELETE")

	// Patient routes
	r.HandleFunc("/patient", handlers.CreatePatient).Methods("POST")
	r.HandleFunc("/patient", handlers.GetAllPatients).Methods("GET")
	r.HandleFunc("/patient/name/{name}", handlers.GetPatientByName).Methods("GET")
	r.HandleFunc("/patient/{id}", handlers.GetPatientByID).Methods("GET")
	r.HandleFunc("/patient/{id}", handlers.UpdatePatient).Methods("PUT")
	r.HandleFunc("/patient/{id}", handlers.DeletePatient).Methods("DELETE")

	// Procedure routes
	r.HandleFunc("/procedure", handlers.CreateProcedure).Methods("POST")
	r.HandleFunc("/procedure", handlers.GetAllProcedures).Methods("GET")
	r.HandleFunc("/procedure/name/{name}", handlers.GetProcedureByName).Methods("GET")
	r.HandleFunc("/procedure/{id}", handlers.GetProcedureByID).Methods("GET")
	r.HandleFunc("/procedure/{id}", handlers.UpdateProcedure).Methods("PUT")
	r.HandleFunc("/procedure/{id}", handlers.DeleteProcedure).Methods("DELETE")

	// Appointment routes
	r.HandleFunc("/appointment", handlers.CreateAppointment).Methods("POST")
	r.HandleFunc("/appointment", handlers.GetAllAppointments).Methods("GET")
	r.HandleFunc("/appointment/{id}", handlers.GetAppointmentByID).Methods("GET")
	r.HandleFunc("/appointment/{id}", handlers.UpdateAppointment).Methods("PUT")
	r.HandleFunc("/appointment/{id}", handlers.DeleteAppointment).Methods("DELETE")

	return r
}
