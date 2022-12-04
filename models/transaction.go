package models

import "time"

type Transaction struct {
	ID        int          `json:"id" gorm:"primary_key:auto_increment"`
	Name      string       `json:"name" form:"name" gorm:"type:text"`
	Email     string       `json:"email" form:"email" gorm:"type: text"`
	Phone     string       `json:"phone" form:"phone" gorm:"type:text"`
	Poscode   string       `json:"poscode" form:"poscode" gorm:"type:text"`
	Address   string       `json:"address" form:"address" gorm:"type:text"`
	Subtotal  int          `json:"price" gorm:"-"`
	Status    string       `json:"status" gorm:"type:varchar(25)"`
	UserID    int          `json:"user_id"`
	User      UserResponse `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	OrderID   int          `json:"order_id"`
	Order     []Order      `json:"order" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
}

type TransactionResponse struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
