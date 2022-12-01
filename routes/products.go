package routes

import (
	"backend/handlers"
	"backend/pkg/middleware"
	"backend/pkg/mysql"
	"backend/repositories"

	"github.com/gorilla/mux"
)

func ProductRoutes(r *mux.Router) {
	ProductRepository := repositories.RepositoryProduct(mysql.DB)
	h := handlers.HandlerProduct(ProductRepository)

	r.HandleFunc("/products", h.FindProducts).Methods("GET")
	r.HandleFunc("/product/{id}", middleware.Auth(h.GetProduct)).Methods("GET")
	r.HandleFunc("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct))).Methods("POST")
	r.HandleFunc("/product/{id}", middleware.Auth(middleware.UploadFile(h.UpdateProduct))).Methods("PATCH")
	r.HandleFunc("/product/{id}", middleware.Auth(h.DeleteProduct)).Methods("DELETE")
}

// {
// 	"products": [
// 	  {
// 		"id": 24,
// 		"qty": 2,
// 		"toppings": [1, 2]
// 	  },
// 	  {
// 		"id": 41,
// 		"qty": 1,
// 		"toppings": [3, 4]
// 	  }
// 	]
// }
