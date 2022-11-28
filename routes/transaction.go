package routes

import (
	"backend/handlers"
	"backend/pkg/mysql"
	"backend/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	TransactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(TransactionRepository)

	r.HandleFunc("/transactions", h.Checkout).Methods("GET")

}
