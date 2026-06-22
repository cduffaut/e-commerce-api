package repository

import (
	"context"
	"fmt"

	"github.com/cduffaut/e-commerce-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository interface {
	AddItem(ctx context.Context, item *domain.CartItem) error
	GetItems(ctx context.Context, userID int64) ([]domain.CartItem, error)
	UpdateItem(ctx context.Context, id int64, userID int64, quantity int) error
	RemoveItem(ctx context.Context, id int64, userID int64) error
	ClearCart(ctx context.Context, userID int64) error
}

type cartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	query := `
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + $3, updated_at = NOW()
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(ctx, query, item.UserID, item.ProductID, item.Quantity).Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *cartRepository) GetItems(ctx context.Context, userID int64) ([]domain.CartItem, error) {
	query := `
		SELECT id, user_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE user_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	defer rows.Close()

	items := []domain.CartItem{}
	for rows.Next() {
		var item domain.CartItem
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID,
			&item.Quantity, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *cartRepository) UpdateItem(ctx context.Context, id int64, userID int64, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity=$1, updated_at=NOW()
		WHERE id=$2 AND user_id=$3
	`

	result, err := r.db.Exec(ctx, query, quantity, id, userID)

	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

func (r *cartRepository) RemoveItem(ctx context.Context, id int64, userID int64) error {
	query := `
		DELETE FROM cart_items WHERE id=$1 AND user_id=$2
	`

	result, err := r.db.Exec(ctx, query, id, userID)

	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("item not found")
	}

	return nil
}

func (r *cartRepository) ClearCart(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM cart_items WHERE user_id=$1`, userID)

	if err != nil {
		return fmt.Errorf("fail to delete item: %w", err)
	}

	return nil
}
