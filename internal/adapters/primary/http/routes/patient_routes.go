package routes

import (
	"go-dentist-api/internal/adapters/primary/http/handlers"

	"github.com/gorilla/mux"
)

func RegisterPatientRoutes(router *mux.Router, patientHandler *handlers.PatientHandler) {
	// Rotas principais de pacientes
	patientRouter := router.PathPrefix("/patients").Subrouter()

	// CRUD básico de pacientes
	patientRouter.HandleFunc("", patientHandler.CreatePatient).Methods("POST")
	patientRouter.HandleFunc("/{id}", patientHandler.GetPatient).Methods("GET")
	patientRouter.HandleFunc("/{id}", patientHandler.UpdatePatient).Methods("PUT")
	patientRouter.HandleFunc("/{id}", patientHandler.DeletePatient).Methods("DELETE")
	patientRouter.HandleFunc("/{id}/deactivate", patientHandler.DeactivatePatient).Methods("PATCH")

	// Busca por CPF
	patientRouter.HandleFunc("/cpf/{cpf}", patientHandler.GetPatientByCPF).Methods("GET")

	// Procedimentos do paciente
	patientRouter.HandleFunc("/{id}/procedures", patientHandler.AddProcedure).Methods("POST")
	patientRouter.HandleFunc("/{id}/procedures", patientHandler.GetPatientProcedures).Methods("GET")

	// Exames do paciente
	patientRouter.HandleFunc("/{id}/exams", patientHandler.AddExam).Methods("POST")
	patientRouter.HandleFunc("/{id}/exams", patientHandler.GetPatientExams).Methods("GET")

	// Histórico médico
	patientRouter.HandleFunc("/{id}/medical-history", patientHandler.UpdateMedicalHistory).Methods("PUT")
	patientRouter.HandleFunc("/{id}/medical-history", patientHandler.GetMedicalHistory).Methods("GET")

	// Rotas específicas de exames (fora do contexto de paciente)
	examRouter := router.PathPrefix("/exams").Subrouter()
	examRouter.HandleFunc("/{examId}", patientHandler.UpdateExam).Methods("PUT")

	// Rotas de clínica para pacientes
	clinicRouter := router.PathPrefix("/clinics").Subrouter()
	clinicRouter.HandleFunc("/{clinicId}/patients", patientHandler.GetPatientsByClinic).Methods("GET")
	clinicRouter.HandleFunc("/{clinicId}/patients/search", patientHandler.SearchPatients).Methods("GET")
	clinicRouter.HandleFunc("/{clinicId}/patients/statistics", patientHandler.GetPatientStatistics).Methods("GET")
}