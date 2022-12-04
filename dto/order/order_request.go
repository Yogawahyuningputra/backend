package orderdto

type OrderRequest struct {
	UserID    int   `json:"user_id" form:"user_id"`
	ProductID int   `json:"product_id" form:"product_id"`
	ToppingID []int `json:"topping_id" form:"topping_id"`
	// TransactionID int   `json:"transaction_id" form:"transaction_id"`
	Subtotal int `json:"subtotal" form:"Subtotal"`
	Qty      int `json:"qty" form:"qty"`
}
type OrderResponse struct {
	// ID        int    `json:"id"`
	// UserID    int    `json:"user_id"`
	// ProductID int    `json:"product_id" form:"product" validate:"required"`
	// ToppingID []int  `json:"topping_id" form:"topping" validate:"required"`
	// Subtotal  int    `json:"subtotal" form:"Subtotal" validate:"required"`
	// Status    string `json:"status"`
	Qty int `json:"qty" form:"qty"`
}
