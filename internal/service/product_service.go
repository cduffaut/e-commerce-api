package service

import (
	"context"
	"fmt"

	"github.com/cduffaut/e-commerce-api/internal/domain"
	"github.com/cduffaut/e-commerce-api/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error)
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	List(ctx context.Context, filter domain.ProductsFilter) ([]domain.Product, error)
	Update(ctx context.Context, id int64, req domain.UpdateProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, id int64) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Category:    req.Category,
		IsActive:    true,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return product, nil
}

func (s *productService) List(ctx context.Context, filter domain.ProductsFilter) ([]domain.Product, error) {

	// if user does not send anything
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.PageSize <= 0 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	return s.repo.List(ctx, filter)
}

func (s *productService) Update(ctx context.Context, id int64, req domain.UpdateProductRequest) (*domain.Product, error) {
	// get the product
	product, err := s.repo.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// only apply changes for non null data given
	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Description != nil {
		product.Description = *req.Description
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.Category != nil {
		product.Category = *req.Category
	}

	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	// check if the product does exist
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}
	return s.repo.Delete(ctx, id)
}
