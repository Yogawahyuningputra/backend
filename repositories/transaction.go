package repositories

import (
	"backend/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransactionRepository interface {
	AddTransaction(transaction models.Transaction) (models.Transaction, error)
	FindTransactions(ID int) ([]models.Transaction, error)
	FindTransactionID(ID int) ([]models.Transaction, error)
	GetTransaction(ID int) (models.Transaction, error)
	GetOrderByUser(ID int) ([]models.Order, error)
	GetOrderByID() ([]models.Transaction, error)
	UpdateTransaction(transaction models.Transaction) (models.Transaction, error)
	CancelTransaction(transaction models.Transaction) (models.Transaction, error)
}

func RepositoryTransaction(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) AddTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Create(&transaction).Error
	return transaction, err
}

func (r *repository) FindTransactions(ID int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.db.Preload("Product").Preload("Topping").Preload("User").Where("UserID = ?", ID).Find(&transactions).Error
	return transactions, err
}

func (r *repository) FindTransactionID(ID int) ([]models.Transaction, error) {
	var transaction []models.Transaction
	err := r.db.Preload("User").Preload(clause.Associations).Preload("Order.Product").Preload("Order.Topping").Where("status != ? AND account_id = ?", "waiting", ID).Find(&transaction).Error

	return transaction, err
}

func (r *repository) GetTransaction(ID int) (models.Transaction, error) {
	var transaction models.Transaction
	err := r.db.Preload("Order").Preload("User").First(&transaction, ID).Error
	return transaction, err
}

func (r *repository) GetOrderByUser(ID int) ([]models.Order, error) {
	var order []models.Order
	err := r.db.Preload("Product").Preload("Topping").Preload("User").Where("user_id =?", ID).Find(&order).Error

	return order, err
}

func (r *repository) GetOrderByID() ([]models.Transaction, error) {
	var order []models.Transaction
	err := r.db.Preload("Order").Preload("User").Where("status != ?", "waiting").Find(&order).Error
	return order, err
}

func (r *repository) UpdateTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Save(&transaction).Error
	return transaction, err
}

func (r *repository) CancelTransaction(transaction models.Transaction) (models.Transaction, error) {
	err := r.db.Delete(&transaction).Error
	return transaction, err
}
