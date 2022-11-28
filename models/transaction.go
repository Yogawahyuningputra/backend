package models

import "time"

type Transaction struct {
	ID        int             `json:"id" gorm:"primary_key:auto_increment"`
	UserID    int             `json:"user_id"`
	User      UserResponse    `json:"user"`
	OrderID   int             `json:"order_id"`
	Order     []OrderResponse `json:"order" gorm:"-"`
	Subtotal  int             `json:"subtotal"`
	Status    string          `json:"status" gorm:"type:varchar(25)"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
}

type TransactionResponse struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

func (TransactionResponse) TableName() string {
	return "transactions"
}
