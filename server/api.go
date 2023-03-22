package server

import (
	"github.com/gorilla/mux"
)

func (s *APIServer) initApiRoutes(r *mux.Router) {
	router := r.PathPrefix("/api").Subrouter()
	router.Methods("GET").Path("/books").Handler(Endpoint{s.ListBooks})
	router.Methods("GET").Path("/books/{id}").Handler(Endpoint{s.GetBookById})
	router.Methods("POST").Path("/books").Handler(Endpoint{s.CreateBook})
	router.Methods("DELETE").Path("/books/{id}").Handler(Endpoint{s.DeleteBook})

	//user endpoint
	router.Methods("POST").Path("/user").Handler(Endpoint{s.CreateUser})
	router.Methods("POST").Path("/login").Handler(Endpoint{s.Login})
	router.Methods("POST").Path("/logout").Handler(Endpoint{s.Logout})
}
