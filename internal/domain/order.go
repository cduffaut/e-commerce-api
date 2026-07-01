package domain

import "time"

type Order struct {
	ID              int64     `db:"id"`
	UserID          int64     `db:"user_id"`
	Status          string    `db:"status"`
	TotalAmount     float64   `db:"total_amount"`
	StripePaymentID string    `db:"stripe_payment_id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

type OrderItem struct {
	ID        int64     `db:"id"`
	OrderID   int64     `db:"order_id"`
	ProductID int64     `db:"product_id"`
	Quantity  int       `db:"quantity"`
	UnitPrice float64   `db:"unit_price"`
	CreatedAt time.Time `db:"created_at"`
}

// enrich OrderItem w/product name
type OrderItemResponse struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

// complete orders with its lines
type OrderResponse struct {
	ID              int64               `json:"id"`
	Status          string              `json:"status"`
	Items           []OrderItemResponse `json:"items"`
	TotalAmount     float64             `json:"total_amount"`
	StripePaymentID string              `json:"stripe_payment_id"`
	CreatedAt       time.Time           `json:"created_at"`
}

// return after order creation -- finalise payment via stripe.js
type CheckoutResponse struct {
	OrderID      int64  `json:"order_id"`
	ClientSecret string `json:"client_secret"`
}
