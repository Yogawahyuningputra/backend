package models

import "time"

type Order struct {
	ID            int                 `json:"id" gorm:"primary_key:auto_increment"`
	Qty           int                 `json:"qty" gorm:"type:int"`
	Subtotal      int                 `json:"subtotal" gorm:"type: int"`
	ProductID     int                 `json:"product_id"`
	Product       ProductResponse     `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Topping       []Topping           `json:"toppings" gorm:"many2many:order_toppings; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID        int                 `json:"user_id"`
	User          UserResponse        `json:"user"`
	TransactionID int                 `json:"transaction_id"`
	Transaction   TransactionResponse `json:"transaction" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt     time.Time           `json:"-"`
	UpdatedAt     time.Time           `json:"-"`
}

type OrderResponse struct {
	ID        int             `json:"id"`
	ProductID int             `json:"product_id"`
	Product   ProductResponse `json:"product"`
	ToppingID []int           `json:"topping_id"`
	Topping   []Topping       `json:"toppings"`
	UserID    int             `json:"user_id"`
	Qty       int             `json:"qty"`
	Subtotal  int             `json:"subtotal"`
}

func (OrderResponse) TableName() string {
	return "orders"
}
