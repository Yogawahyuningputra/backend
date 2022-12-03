package models

import "time"

type Order struct {
	ID        int             `json:"id" gorm:"primary_key:auto_increment"`
	Qty       int             `json:"qty"`
	Subtotal  int             `json:"subtotal" form:"subtotal" gorm:"type: int"`
	ProductID int             `json:"product_id"`
	Product   ProductResponse `json:"product" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// ToppingID []int           `json:"topping_id" form:"topping_id" gorm:"-"`
	Topping   []Topping    `json:"toppings" gorm:"many2many:order_toppings; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    int          `json:"user_id"`
	User      UserResponse `json:"user"`
	Price     int          `json:"price" gorm:"type:int"`
	CreatedAt time.Time    `json:"-"`
	UpdatedAt time.Time    `json:"-"`
}

type OrderResponse struct {
	ID        int             `json:"id"`
	ProductID int             `json:"product_id"`
	Product   ProductResponse `json:"product"`
	ToppingID []int           `json:"topping_id" form:"topping_id" gorm:"-"`
	Topping   []Topping       `json:"toppings" gorm:"many2many:order_toppings; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    int             `json:"user_id"`
	Price     int             `json:"price" gorm:"type:int"`
	Qty       int             `json:"qty"`
	Subtotal  int             `json:"subtotal"`
}

func (OrderResponse) TableName() string {
	return "orders"
}
