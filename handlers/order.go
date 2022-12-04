package handlers

import (
	orderdto "backend/dto/order"
	resultdto "backend/dto/result"
	"backend/models"
	"backend/repositories"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type handlerOrder struct {
	OrderRepository repositories.OrderRepository
}

func HandlerOrder(OrderRepository repositories.OrderRepository) *handlerOrder {
	return &handlerOrder{OrderRepository}
}

func (h *handlerOrder) FindOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	orders, err := h.OrderRepository.FindOrders()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: orders}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerOrder) GetOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	cart, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: (cart)}
	json.NewEncoder(w).Encode(response)
}
func (h *handlerOrder) GetOrderById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// id, _ := strconv.Atoi(mux.Vars(r)["id"])
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int64(userInfo["id"].(float64))

	transaction, _ := h.OrderRepository.GetTransactionID(int(userID))

	cart, err := h.OrderRepository.GetOrderById(int(transaction.ID))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: (cart)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerOrder) CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	request := new(orderdto.OrderRequest)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "cek dto"}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err := validation.Struct(request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "error validation"}
		json.NewEncoder(w).Encode(response)
		return
	}

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	product, err := h.OrderRepository.GetProductOrder(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Product Not Found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	toppings, err := h.OrderRepository.GetToppingOrder(request.ToppingID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Topping Not Found!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var priceTopings = 0
	for _, i := range toppings {
		priceTopings += i.Price
	}
	var subTotal = request.Qty * (product.Price + priceTopings)

	CekRequestTrans, _ := h.OrderRepository.GetTransactionID(userID)

	var transID int
	if CekRequestTrans.ID != 0 {
		transID = CekRequestTrans.ID
	} else {
		requestTrans := models.Transaction{
			Name:     "-",
			Email:    "-",
			Phone:    "-",
			Poscode:  "-",
			Address:  "-",
			Status:   "Waiting",
			Subtotal: 0,
			UserID:   userID,
		}
		transOrder, _ := h.OrderRepository.RequestTransaction(requestTrans)
		transID = transOrder.ID
	}
	// fmt.Println("cek = ", transID)
	dataOrder := models.Order{
		UserID:        userID,
		ProductID:     product.ID,
		Topping:       toppings,
		Qty:           request.Qty,
		Subtotal:      subTotal,
		TransactionID: transID,
	}

	cart, err := h.OrderRepository.CreateOrder(dataOrder)

	fmt.Println("Cart = ", cart.Qty)
	fmt.Println("Subtotal = ", cart.Subtotal)
	fmt.Println("TransactionID = ", cart.TransactionID)
	fmt.Println("ProductID = ", cart.ProductID)
	fmt.Println("Topping = ", cart.Topping)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: "Order Failed!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	order, _ := h.OrderRepository.GetOrder(cart.ID)

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: order}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerOrder) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	order, err := h.OrderRepository.GetOrder(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.OrderRepository.DeleteOrder(order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: convertResponseOrder(data)}
	json.NewEncoder(w).Encode(response)
}

func convertResponseOrder(u models.Order) models.Order {
	return models.Order{
		ID:       u.ID,
		Qty:      u.Qty,
		Subtotal: u.Subtotal,
		Product:  u.Product,
		Topping:  u.Topping,
	}
}
