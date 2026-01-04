package services

import (
	"context"
	"products-api/internal/models"
	"products-api/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	repo *repository.ProductRepository
}

func (s *ProductService) Create(context context.Context, product models.Product) error {
	_, err := s.repo.Create(context, product)
	return err
}

// GetProducts retrieves all products
func (s *ProductService) GetProducts(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductService) UpdateProductCount(ctx context.Context, id string, sold int) (*models.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.repo.UpdateProductCount(ctx, product, sold)
}

// CreateProduct creates a new product
// func (s *ProductService) CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error) {
// 	return s.repo.CreateProduct(ctx, req)
// }
