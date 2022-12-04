package routes

import (
	"backend/handlers"
	"backend/pkg/middleware"
	"backend/pkg/mysql"
	"backend/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	TransactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(TransactionRepository)

	// r.HandleFunc("/transaction", middleware.Auth(h.AddTransaction)).Methods("POST")
	// r.HandleFunc("/transaction/{id}", middleware.Auth(h.CancelTransaction)).Methods("PATCH")
	// r.HandleFunc("/transaction/{id}", middleware.Auth(h.UpdateTransaction)).Methods("PATCH")
	// r.HandleFunc("/transaction/{id}", middleware.Auth(h.AcceptTransaction)).Methods("PATCH")

	r.HandleFunc("/transactions", middleware.Auth(h.FindTransactions)).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.AddTransaction)).Methods("POST")
	r.HandleFunc("/canceltrans/{id}", middleware.Auth(h.CancelTransaction)).Methods("PATCH")
	r.HandleFunc("/updatetrans/{id}", middleware.Auth(h.UpdateTransaction)).Methods("PATCH")
	r.HandleFunc("/accepttrans/{id}", middleware.Auth(h.AcceptTransaction)).Methods("PATCH")
	r.HandleFunc("/my-order", middleware.Auth(h.GetOrderByID)).Methods("GET")
	// // }

}

// func TransactionRoutes(r *mux.Router) {
// 	TransactionRepository := repositories.RepositoryTransaction(mysql.DB)
// 	h := handlers.HandlerTransaction(TransactionRepository)
//
