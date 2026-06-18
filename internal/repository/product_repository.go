package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/cduffaut/e-commerce-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context, filter []domain.Product) ([]domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
	UpdateStock(ctx context.Context, id int64, delta int) error
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, category, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.Stock, product.Category, product.IsActive).Scan(
		&product.ID, &product.CreatedAt, &product.UpdatedAt)
}

func (r *productRepository) GetbyID(ctx context.Context, id int64) (*domain.Product, err) {
	query := `
		SELECT id, name, description, price, stock, category, is_active, created_at, updated_at
		FROM products
		WHERE id = $1 AND is_active == true
	`
	p := &domain.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return p, nil
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, error) {
	conditions := []string{"is_active = true"}
	args := []any{}
	argIndex := 1

	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%" + filter.Search + "%")
		argIndex++
	}

	if filter.Category != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%" + filter.Category + "%")
		argIndex++
	}

	// pagination
	offset := (filter.Page - 1) * filter.PageSize
	args = append(args, filter.PageSize, offset)

	query := fmt.Sprintf(`
		SELECT id, name, description, price, stock, category, is_active, created_at, updated_at
		FROM products
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, strings.Join(conditions, " AND "), argIndex, argIndex+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		var p domain.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.Category, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name=$1, description=$2, price=$3, stock=$4, category=$5, is_active=$6, updated_at=NOW()
		WHERE id=$7
		RETURNING updated_at
	`
	return r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price,
		product.Stock, product.Category, product.IsActive, product.ID
		).Scan(&product.UpdatedAt)
}

// delete (deactivate) product to preserve stock tracking history
func (r *productRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "UPDATE products SET is_active=false, updated_at=NOW() WHERE id=$1", id)

	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id int64, delta int) error {
	query := `
		UPDATE products
		SET stock = stock + $1, updated_at = NOW()
		WHERE id = $2 AND stock + $1 >= 0
	`

	result, err := r.db.Exec(ctx, query, delta, id)

	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	// delta can be positive or negative if delta conduce to
	// negative result, it protect from negative stock
	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock for product %d", id)
	}

	return nil
}
