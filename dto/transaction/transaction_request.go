package transactiondto

type Transaction struct {
	Name     string `json:"name" form:"name"`
	Address  string `json:"address" form:"address"`
	OrderID  []int  `json:"order_id" form:"order_id"`
	Subtotal string `json:"subtotal" form:"subtotal"`
	Status   string `json:"status" form:"status"`
	UserID   int    `json:"user_id" form:"user_id"`
}
