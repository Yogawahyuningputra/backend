package transactiondto

type Transaction struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Poscode  string `json:"poscode" form:"poscode"`
	Address  string `json:"address" form:"address"`
	OrderID  []int  `json:"order_id" form:"order_id"`
	Subtotal string `json:"subtotal" form:"subtotal"`
	Status   string `json:"status" form:"status"`
	UserID   int    `json:"user_id" form:"user_id"`
}

type UpdateTransaction struct {
	Name    string `json:"name" form:"name"`
	Email   string `json:"email" form:"email"`
	Phone   string `json:"phone" form:"phone"`
	Poscode string `json:"poscode" form:"poscode"`
	Address string `json:"address" form:"address"`
	// Status string `json:"status"`
}
