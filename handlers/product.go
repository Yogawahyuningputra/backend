package handlers

import (
	productdto "backend/dto/product"
	resultdto "backend/dto/result"
	"backend/models"
	"backend/repositories"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type handlerProduct struct {
	ProductRepository repositories.ProductRepository
}

// var path_file = os.Getenv("PATH_FILE")

var path_file = "http://localhost:5000/uploads/"

func HandlerProduct(ProductRepository repositories.ProductRepository) *handlerProduct {
	return &handlerProduct{ProductRepository}
}

func (h *handlerProduct) FindProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	products, err := h.ProductRepository.FindProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	var ProductResponses []productdto.ProductResponse
	for _, s := range products {
		ProductResponse := productdto.ProductResponse{
			ID:    s.ID,
			Title: s.Title,
			Price: s.Price,
			Image: s.Image,
		}

		ProductResponses = append(ProductResponses, ProductResponse)
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ProductResponses}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerProduct) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	products, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
	}

	ProductResponse := productdto.ProductResponse{
		ID:    products.ID,
		Title: products.Title,
		Price: products.Price,
		Image: products.Image,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ProductResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerProduct) CreateProduct(w http.ResponseWriter, r *http.Request) {
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
	filepath := dataContex.(string)

	price, _ := strconv.Atoi(r.FormValue("price"))
	request := productdto.ProductRequest{
		Title: r.FormValue("title"),
		Price: price,
		Image: filepath,
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	var ctx = context.Background()
	// var CLOUD_NAME = os.Getenv("CLOUD_NAME")
	// var API_KEY = os.Getenv("API_KEY")
	// var API_SECRET = os.Getenv("API_SECRET")

	// Add your Cloudinary credentials ...
	cld, _ := cloudinary.NewFromParams("doj9z1rop", "533982769458514", "Jdekzp9W3K0tT_xVUjR8BS0c5xA")
	// cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)

	// Upload file to Cloudinary ...
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "waysbucks/product"})

	if err != nil {
		fmt.Println(err.Error())
	}

	product := models.Product{
		Title: request.Title,
		Price: request.Price,
		Image: resp.SecureURL,
		// UserID:      userId,

	}

	product, err = h.ProductRepository.CreateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ProductResponse := productdto.ProductResponse{
		ID:    product.ID,
		Title: product.Title,
		Price: product.Price,
		Image: path_file + product.Image,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ProductResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerProduct) UpdateProduct(w http.ResponseWriter, r *http.Request) {
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

	request := productdto.ProductRequest{
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

	product, err := h.ProductRepository.GetProduct(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Title != "" {
		product.Title = request.Title
	}

	if request.Price != 0 {
		product.Price = request.Price
	}

	if request.Image != "empty" {
		product.Image = request.Image
	}

	data, err := h.ProductRepository.UpdateProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ProductResponse := productdto.ProductResponse{
		ID:    data.ID,
		Title: data.Title,
		Price: data.Price,
		Image: path_file + data.Image,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ProductResponse}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerProduct) DeleteProduct(w http.ResponseWriter, r *http.Request) {
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

	product, err := h.ProductRepository.GetProduct(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := resultdto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.ProductRepository.DeleteProduct(product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := resultdto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	ProductDeleteResponse := productdto.ProductDeleteResponse{
		ID:    data.ID,
		Title: data.Title,
	}

	w.WriteHeader(http.StatusOK)
	response := resultdto.SuccessResult{Code: http.StatusOK, Data: ProductDeleteResponse}
	json.NewEncoder(w).Encode(response)

}
