package repository

import (
	"context"
	"fmt"

	"github.com/cduffaut/e-commerce-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error
	GetByID(ctx context.Context, id int64, userID int64) (*domain.Order, error)
	GetItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error)
	ListByUser(ctx context.Context, userID int64) ([]domain.Order, error)
	UpdateStatus(ctx context.Context, id int64, status string, stripePaymentID string) error
}

type orderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *domain.Order, items []domain.OrderItem) error {
	tx, err := r.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, `INSERT INTO orders (user_id, status, total_amount, stripe_payment_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`, order.UserID, order.Status, order.TotalAmount, order.StripePaymentID).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	for i := range items {
		items[i].OrderID = order.ID

		err = tx.QueryRow(ctx, `
			INSERT INTO order_items (order_id, product_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at
		`, items[i].OrderID, items[i].ProductID, items[i].Quantity, items[i].UnitPrice).Scan(&items[i].ID, &items[i].CreatedAt)

		if err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *orderRepository) GetByID(ctx context.Context, id int64, userID int64) (*domain.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, stripe_payment_id, created_at, updated_at
		FROM orders
		WHERE id=$1 AND user_id=$2
	`

	o := &domain.Order{}
	err := r.db.QueryRow(ctx, query, id, userID).Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.StripePaymentID, &o.CreatedAt, &o.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}
	return o, nil
}

func (r *orderRepository) GetItems(ctx context.Context, orderID int64) ([]domain.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, unit_price, created_at
		FROM order_items
		WHERE order_id=$1
	`

	rows, err := r.db.Query(ctx, query, orderID)

	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	defer rows.Close()

	items := []domain.OrderItem{}
	for rows.Next() {
		var item domain.OrderItem

		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice, &item.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order item: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *orderRepository) ListByUser(ctx context.Context, userID int64) ([]domain.Order, error) {
	query := `
		SELECT id, user_id, status, total_amount, stripe_payment_id, created_at, updated_at
		FROM orders
		WHERE user_id=$1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	defer rows.Close()

	orders := []domain.Order{}

	for rows.Next() {
		var o domain.Order

		err := rows.Scan(&o.ID, &o.UserID, &o.Status, &o.TotalAmount, &o.StripePaymentID, &o.CreatedAt, &o.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}

		orders = append(orders, o)
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, id int64, status string, stripePaymentID string) error {
	_, err := r.db.Exec(ctx, `
		UPDATE orders
		SET status=$1, stripe_payment_id=$2, updated_at=NOW()
		WHERE id=$3
	`, status, stripePaymentID, id)

	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	return nil
}
