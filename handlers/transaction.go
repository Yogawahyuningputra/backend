package handlers

import (
	resultdto "backend/dto/result"
	transactiondto "backend/dto/transaction"
	"backend/models"
	"backend/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}
func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transactions, err := h.TransactionRepository.GetOrderByID()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusBadRequest, Data: transactions}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) AddTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	request := new(transactiondto.Transaction)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validate := validator.New()
	err := validate.Struct(request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Validation Error"}
		json.NewEncoder(w).Encode(response)
		return
	}
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	userCart, _ := h.TransactionRepository.GetOrderByUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Purchased Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}
	var totalTransaction = 0

	for _, i := range userCart {
		totalTransaction += i.Subtotal
	}

	dataOrders := models.Transaction{
		UserID:   userID,
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Poscode:  request.Poscode,
		Address:  request.Address,
		Order:    userCart,
		Subtotal: totalTransaction,
		Status:   "Waiting",
	}

	itemTransaction, err := h.TransactionRepository.AddTransaction(dataOrders)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	orderItem, _ := h.TransactionRepository.GetTransaction(itemTransaction.ID)

	data := models.Transaction{
		ID:       orderItem.ID,
		Name:     request.Name,
		Email:    request.Email,
		Phone:    request.Phone,
		Poscode:  request.Poscode,
		Address:  request.Address,
		Order:    userCart,
		Subtotal: totalTransaction,
		// Status:   "pending",
		// UserID:   orderItem.UserID,
		// User:     orderItem.User,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)

}
func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	request := new(transactiondto.UpdateTransaction)
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	transaction, err := h.TransactionRepository.GetTransaction(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// if request.Status != "" {
	// 	transaction.Status = request.Status
	// }
	orderTrans, err := h.TransactionRepository.GetOrderByUser(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}

	var Total = 0
	for _, i := range orderTrans {
		Total += i.Subtotal
	}

	// fmt.Println(orderTrans)
	// fmt.Println(transaction.ID)
	// fmt.Println(Total)

	if request.Name != "" {
		transaction.Name = request.Name
	}
	if request.Email != "" {
		transaction.Email = request.Email
	}
	if request.Phone != "" {
		transaction.Phone = request.Phone
	}
	if request.Poscode != "" {
		transaction.Poscode = request.Poscode
	}
	if request.Address != "" {
		transaction.Address = request.Address
	}
	transaction.Status = "Payment"
	transaction.Subtotal = Total

	data, err := h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	trans, _ := h.TransactionRepository.GetTransaction(data.ID)
	orderUser, err := h.TransactionRepository.GetOrderByUser(data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}
	dataUpdate := models.Transaction{
		ID:       trans.ID,
		Name:     trans.Name,
		Address:  trans.Address,
		Status:   trans.Status,
		Order:    orderUser,
		Subtotal: trans.Subtotal,
		UserID:   trans.UserID,
		User:     trans.User,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: dataUpdate}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) AcceptTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Check ID Transaction"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// fmt.Println(transaction.ID)
	// fmt.Println(transaction.Status)

	transaction.Status = "Success"
	data, err := h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CancelTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transaction.Status = "Cancel"
	data, err := h.TransactionRepository.CancelTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}
func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	trans, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "Failed"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// order, err := h.TransactionRepository.GetOrderByID(trans.UserID)

	dataTransactions := models.Transaction{
		ID:      trans.ID,
		Name:    trans.Name,
		Email:   trans.Email,
		Phone:   trans.Phone,
		Address: trans.Address,

		Subtotal: trans.Subtotal,
		UserID:   trans.UserID,
		Status:   trans.Status,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: dataTransactions}
	json.NewEncoder(w).Encode(response)
}
func (h *handlerTransaction) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	orders, err := h.TransactionRepository.FindTransactionID(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusAccepted, Data: orders}
	json.NewEncoder(w).Encode(response)
}
