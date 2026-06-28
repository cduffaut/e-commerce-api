package service

import (
	"context"
	"fmt"

	"github.com/cduffaut/e-commerce-api/internal/domain"
	"github.com/cduffaut/e-commerce-api/internal/repository"
)

type CartService interface {
	AddItem(ctx context.Context, userID int64, req domain.AddToCartRequest) error
	GetCart(ctx context.Context, userID int64) (*domain.CartResponse, error)
	UpdateItem(ctx context.Context, userID int64, itemID int64, req domain.UpdateCartItemRequest) error
	RemoveItem(ctx context.Context, userID int64, itemID int64) error
}

type cartService struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository) CartService {
	return &cartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *cartService) AddItem(ctx context.Context, userID int64, req domain.AddToCartRequest) error {
	// check if product exist and is available
	product, err := s.productRepo.GetByID(ctx, req.ProductID)

	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// enough stock
	if product.Stock < req.Quantity {
		return fmt.Errorf("insufficient stock: requested %d, available %d", req.Quantity, product.Stock)
	}

	item := &domain.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	return s.cartRepo.AddItem(ctx, item)
}

func (s *cartService) GetCart(ctx context.Context, userID int64) (*domain.CartResponse, error) {
	items, err := s.cartRepo.GetItems(ctx, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	cartItems := []domain.CartItemResponse{}

	var totalPrice float64

	for _, item := range items {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)

		if err != nil {
			continue
		}

		cartItems = append(cartItems, domain.CartItemResponse{
			ID:       item.ID,
			Product:  *product,
			Quantity: item.Quantity,
		})

		totalPrice += product.Price * float64(item.Quantity)
	}

	return &domain.CartResponse{
		Items:      cartItems,
		TotalPrice: totalPrice,
	}, nil
}

func (s *cartService) UpdateItem(ctx context.Context, userID int64, itemID int64, req domain.UpdateCartItemRequest) error {
	return s.cartRepo.UpdateItem(ctx, itemID, userID, req.Quantity)
}

func (s *cartService) RemoveItem(ctx context.Context, userID int64, itemID int64) error {
	return s.cartRepo.RemoveItem(ctx, itemID, userID)
}
