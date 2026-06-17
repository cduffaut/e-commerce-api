package domain

import "time"

type CartItem struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	ProductID int64     `db:"product_id"`
	Quantity  int       `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// product info (enrich CartItem)
type CartItemResponse struct {
	ID       int64   `json:"id"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

// user complete cart
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"total_price"`
}

type AddToCartRequest struct {
	ProductID int64 `json:"product_id" binding:"required,gt=0"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}
