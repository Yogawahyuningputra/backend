package routes

import (
	"backend/handlers"
	"backend/pkg/middleware"
	"backend/pkg/mysql"
	"backend/repositories"

	"github.com/gorilla/mux"
)

func OrderRoutes(r *mux.Router) {
	OrderRepository := repositories.RepositoryOrder(mysql.DB)
	h := handlers.HandlerOrder(OrderRepository)

	r.HandleFunc("/orders", middleware.Auth(h.FindOrders)).Methods("GET")
	r.HandleFunc("/order/{id}", middleware.Auth(h.GetOrder)).Methods("GET")
	r.HandleFunc("/orders-id", middleware.Auth(h.GetOrderById)).Methods("GET")

	r.HandleFunc("/order/{id}", middleware.Auth(h.CreateOrder)).Methods("POST")
	r.HandleFunc("/order/{id}", middleware.Auth(h.DeleteOrder)).Methods("DELETE")

}
