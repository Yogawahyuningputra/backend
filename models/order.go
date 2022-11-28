package models

import "time"

type Order struct {
	ID        int               `json:"id" gorm:"primary_key:auto_increment"`
	Qty       int               `json:"qty"`
	Subtotal  int               `json:"subtotal" form:"subtotal" gorm:"type: int"`
	ProductID int               `json:"product_id"`
	Product   ProductResponse   `json:"product"`
	ToppingID []int             `json:"topping_id" form:"topping_id" gorm:"-"`
	Topping   []ToppingResponse `json:"toppings" gorm:"many2many:order_toppings"`
	UserID    int               `json:"user_id"`
	User      UserResponse      `json:"user"`
	Price     int               `json:"price" gorm:"type:int"`
	// TransactionID int                 `json:"transaction_id"`
	// Transaction   TransactionResponse `json:"transaction"`
	// Status    string    `json:"status"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type OrderResponse struct {
	ID        int               `json:"id"`
	ProductID int               `json:"product_id"`
	Product   ProductResponse   `json:"product"`
	ToppingID []int             `json:"topping_id" form:"topping_id" gorm:"-"`
	Topping   []ToppingResponse `json:"toppings" gorm:"many2many:order_toppings"`
	UserID    int               `json:"user_id"`
	// TransactionID int                 `json:"transaction_id"`
	// Transaction   TransactionResponse `json:"transaction" gorm:"foreignKey:TransactionID"`
	Qty      int `json:"qty"`
	Subtotal int `json:"subtotal"`
}

func (OrderResponse) TableName() string {
	return "orders"
}
