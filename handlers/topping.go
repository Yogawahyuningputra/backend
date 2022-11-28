package handlers

import (
	resultdto "backend/dto/result"
	toppingdto "backend/dto/topping"
	"backend/models"
	"backend/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type handlerTopping struct {
	ToppingRepository repositories.ToppingRepository
}

// var path_topping = os.Getenv("PATH_FILE")
var path_topping = "http://localhost:5000/uploads/"

func HandlerTopping(ToppingRepository repositories.ToppingRepository) *handlerTopping {
	return &handlerTopping{ToppingRepository}
}

func (h *handlerTopping) FindToppings(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	toppings, err := h.ToppingRepository.FindToppings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	var ToppingResponses []toppingdto.ToppingResponse
	for _, s := range toppings {
		ToppingResponse := toppingdto.ToppingResponse{
			ID:    s.ID,
			Title: s.Title,
			Price: s.Price,
			Image: path_topping + s.Image,
		}

		ToppingResponses = append(ToppingResponses, ToppingResponse)
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ToppingResponses}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTopping) GetTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	toppings, err := h.ToppingRepository.GetTopping(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	ProductResponse := toppingdto.ToppingResponse{
		ID:    toppings.ID,
		Title: toppings.Title,
		Price: toppings.Price,
		Image: path_topping + toppings.Image,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ProductResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTopping) CreateTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//get data user token
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)

	userRole := userInfo["role"]

	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := resultdto.ErrorResult{Code: http.StatusUnauthorized, Message: "not an admin"}
		json.NewEncoder(w).Encode(response)
		return
	}
	//get data from middleware
	dataContex := r.Context().Value("dataFile")
	filename := dataContex.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))
	request := toppingdto.ToppingRequest{
		Title: r.FormValue("title"),
		Price: price,
		Image: filename,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	topping := models.Topping{
		Title: request.Title,
		Price: request.Price,
		Image: filename,
	}

	topping, err = h.ToppingRepository.CreateTopping(topping)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ToppingResponse := toppingdto.ToppingResponse{
		ID:    topping.ID,
		Title: topping.Title,
		Price: topping.Price,
		Image: path_topping + topping.Image,
	}

	// product, _ = h.ProductRepository.GetProduct(product.ID)

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ToppingResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTopping) UpdateTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := resultdto.ErrorResult{Code: http.StatusUnauthorized, Message: "not an admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	dataContex := r.Context().Value("dataFile")
	filename := dataContex.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))

	// fmt.Println(filename)
	// fmt.Println(price)

	request := toppingdto.ToppingRequest{
		Title: r.FormValue("title"),
		Price: price,
		Image: filename,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	topping, err := h.ToppingRepository.GetTopping(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Title != "" {
		topping.Title = request.Title
	}

	if request.Price != 0 {
		topping.Price = request.Price
	}

	if request.Image != "empty" {
		topping.Image = request.Image
	}

	data, err := h.ToppingRepository.UpdateTopping(topping)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ToppingResponse := toppingdto.ToppingResponse{
		ID:    data.ID,
		Title: data.Title,
		Price: data.Price,
		Image: path_topping + data.Image,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ToppingResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTopping) DeleteTopping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	if userRole != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		response := resultdto.ErrorResult{Code: http.StatusUnauthorized, Message: "not an admin"}
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	topping, err := h.ToppingRepository.GetTopping(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.ToppingRepository.DeleteTopping(topping)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ToppingDeleteResponse := toppingdto.ToppingDeleteResponse{
		ID:    data.ID,
		Title: data.Title,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ToppingDeleteResponse}
	json.NewEncoder(w).Encode(response)

}
