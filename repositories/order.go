package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	FindOrders() ([]models.Order, error)
	GetOrder(ID int) (models.Order, error)
	GetOrderById(TransactionID int) ([]models.Order, error)
	CreateOrder(order models.Order) (models.Order, error)
	UpdateOrder(order models.Order) (models.Order, error)
	DeleteOrder(order models.Order) (models.Order, error)
	RequestTransaction(transaction models.Transaction) (models.Transaction, error)
	GetProductOrder(ID int) (models.Product, error)
	GetToppingOrder(ID []int) ([]models.Topping, error)
	GetTransactionID(ID int) (models.Transaction, error)
}

func RepositoryOrder(db *gorm.DB) *repository {
	return &repository{db}
}
func (r *repository) FindOrders() ([]models.Order, error) {
	var order []models.Order
	err := r.db.Preload("Product").Preload("Topping").Preload("User").Where("user_id").Find(&order).Error

	return order, err
}
func (r *repository) GetOrder(ID int) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("Product").Preload("Topping").First(&order, ID).Error

	return order, err
}
func (r *repository) GetOrderById(ID int) ([]models.Order, error) {
	var order []models.Order
	err := r.db.Preload("Product").Preload("Topping").Find(&order, "transaction_id = ?", ID).Error

	return order, err
}

func (r *repository) CreateOrder(order models.Order) (models.Order, error) {
	err := r.db.Create(&order).Error

	return order, err
}

func (r *repository) UpdateOrder(order models.Order) (models.Order, error) {
	err := r.db.Save(&order).Error

	return order, err
}
func (r *repository) DeleteOrder(order models.Order) (models.Order, error) {
	err := r.db.Delete(&order).Error

	return order, err
}

func (r *repository) GetProductOrder(ID int) (models.Product, error) {
	var product models.Product
	err := r.db.First(&product, ID).Error
	return product, err
}

func (r *repository) GetToppingOrder(ID []int) ([]models.Topping, error) {
	var topping []models.Topping
	err := r.db.Find(&topping, ID).Error
	return topping, err
}
func (r *repository) GetTransactionID(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("User").Preload("Order").Where("status = ? AND user_id = ?", "waiting", ID).First(&transaction).Error
	return transaction, err
}

func (r *repository) RequestTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error

	return transaction, err
}
