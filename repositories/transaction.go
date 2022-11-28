package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Checkout(transaction *models.Transaction) (*models.Transaction, error)
	GetOrderByUser(ID int) ([]models.Order, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Checkout(transaction *models.Transaction) (*models.Transaction, error) {
	err := r.db.Create(transaction).Error
	return transaction, err
}
func (r *repository) GetOrderByUser(ID int) ([]models.Order, error) {
	var order []models.Order
	err := r.db.Preload("user").Where("user_id =?", ID).Find(&order).Error

	return order, err
}
