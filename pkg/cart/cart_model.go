package cart

import "time"

type CartItemModel struct {
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=1"`
}

type CartModel struct {
	UserID    string          `json:"user_id"`
	Items     []CartItemModel `json:"items"`
	Status    string          `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
}

type NewCartItemModel struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"min=1,required"`
}
