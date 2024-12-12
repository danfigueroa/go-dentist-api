package router

import (
	"go-dentist-api/internal/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/dentist", handlers.CreateDentist).Methods(http.MethodPost)
	r.HandleFunc("/dentists", handlers.GetAllDentists).Methods(http.MethodGet)
	r.HandleFunc("/dentist/{id}", handlers.GetDentistByID).Methods(http.MethodGet)
	r.HandleFunc("/dentist/{id}", handlers.UpdateDentist).Methods(http.MethodPut)
	r.HandleFunc("/dentist/{id}", handlers.DeleteDentist).Methods(http.MethodDelete)

	return r
}
