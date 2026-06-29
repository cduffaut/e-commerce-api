package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"

	"github.com/cduffaut/e-commerce-api/internal/domain"
	"github.com/cduffaut/e-commerce-api/internal/repository"
)

type OrderService interface {
	Checkout(ctx context.Context, userID int64) (*domain.CheckoutResponse, error)
	GetOrder(ctx context.Context, userID int64, orderID int64) (*domain.OrderResponse, error)
	ListOrders(ctx context.Context, userID int64) ([]domain.OrderResponse, error)
	ConfirmPayment(ctx context.Context, stripePaymentID string) error
}

type orderService struct {
	orderRepo   repository.OrderRepository
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	stripeKey   string
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, productRepo repository.ProductRepository,
	stripeKey string) OrderService {
	return &orderService{orderRepo: orderRepo, cartRepo: cartRepo, productRepo: productRepo, stripeKey: stripeKey}
}

func (s *orderService) Checkout(ctx context.Context, userID int64) (*domain.CheckoutResponse, error) {
	// get the cart
	cartItems, err := s.cartRepo.GetItems(ctx, userID)
	if err != nil || len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// build command line + total calculation
	orderItems := []domain.OrderItem{}
	var totalAmount float64

	for _, cartItem := range cartItems {
		product, err := s.productRepo.GetByID(ctx, cartItem.ProductID)

		if err != nil {
			return nil, fmt.Errorf("product %d not found", cartItem.ProductID)
		}

		if product.Stock < cartItem.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s", product.Name)
		}

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: cartItem.ProductID,
			UnitPrice: product.Price,
			Quantity:  cartItem.Quantity,
		})

		totalAmount += product.Price * float64(cartItem.Quantity)
	}
	// payment -- stripe is in cents not euros
	stripe.Key = s.stripeKey
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(totalAmount * 100)),
		Currency: stripe.String(string(stripe.CurrencyEUR)),
	}

	pi, err := paymentintent.New(params)

	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	pi_id, err := strconv.ParseFloat(pi.ID, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert string to float for payment %s", pi.ID)
	}

	// create command in DB
	order := &domain.Order{
		UserID:          userID,
		Status:          "pending",
		TotalAmount:     totalAmount,
		StripePaymentID: pi_id,
	}

	if err := s.orderRepo.Create(ctx, order, orderItems); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// update stock
	for _, item := range orderItems {
		if err := s.productRepo.UpdateStock(ctx, item.OrderID, -item.Quantity); err != nil {
			return nil, fmt.Errorf("failed to update stock: %w", err)
		}
	}

	// clear user cart
	if err := s.cartRepo.ClearCart(ctx, userID); err != nil {
		return nil, fmt.Errorf("failed to clear cart items: %w", err)
	}

	return &domain.CheckoutResponse{OrderID: order.ID, ClientSecret: pi.ClientSecret}, nil
}
