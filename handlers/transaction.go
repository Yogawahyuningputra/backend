package handlers

import (
	resultdto "backend/dto/result"
	transactiondto "backend/dto/transaction"
	"backend/repositories"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) Checkout(w http.ResponseWriter, r *http.Request) {
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
	accountID := int(userInfo["id"].(float64))

	UserCart, err := h.TransactionRepository.GetOrderByUser(accountID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: "User Account not make a Purchase!"}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Println(accountID)
	fmt.Println(UserCart)

}
