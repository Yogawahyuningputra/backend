package routes

import (
	"backend/handlers"
	"backend/pkg/mysql"
	"backend/repositories"

	"github.com/gorilla/mux"
)

func OrderRoutes(r *mux.Router) {
	OrderRepository := repositories.RepositoryOrder(mysql.DB)
	h := handlers.HandlerOrder(OrderRepository)

	r.HandleFunc("/orders", h.FindOrders).Methods("GET")
	r.HandleFunc("/order/{id}", h.GetOrder).Methods("GET")
	r.HandleFunc("/order/{id}", h.CreateOrder).Methods("POST")
	r.HandleFunc("/order/{id}", h.DeleteOrder).Methods("DELETE")

}
