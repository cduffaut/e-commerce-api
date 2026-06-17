package domain

import "time"

type Product struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float64   `db:"price"`
	Stock       int       `db:"stock"`
	Category    string    `db:"category"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// resered to admin
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Stock       int     `json:"stock" binding:"omitempty,gt=0"`
	Category    string  `json:"category"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" binding:"omitempty, gt=0"`
	Stock       *int     `json:"stock" binding:"omitempty, gte=0"`
	Category    *string  `json:"category"`
	IsActive    *bool    `json:"is_active"`
}

// parameters for filters/research
type ProductFilter struct {
	Search   string `form:"q"`
	Category string `form:"category"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}
