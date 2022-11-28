package toppingdto

type ToppingRequest struct {
	Title string `json:"title" form:"title" validate:"required"`
	Price int    `json:"price" form:"price" validate:"required"`
	Image string `json:"image" form:"image" validate:"required"`
}
