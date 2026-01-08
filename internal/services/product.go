package services

import (
	"context"
	"log"
	"products-api/internal/models"
	"products-api/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	repo *repository.ProductRepository
}

func (s *ProductService) Create(context context.Context, product *models.Product) error {
	err := s.repo.Create(context, product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
	}
	return err
}

// GetProducts retrieves all products
func (s *ProductService) GetProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		return nil, err
	}
	log.Printf("Retrieved %d products %v", len(products), products)
	return products, nil
}

func (s *ProductService) UpdateProductCount(ctx context.Context, id string, sold int) (*models.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieving product by ID: %v", err)
		return nil, err
	}
	log.Printf("Updating product count for product ID %s, sold: %d", id, sold)
	err = s.repo.UpdateProductCount(ctx, product, sold)
	if err != nil {
		log.Printf("Error updating product count: %v", err)
		return nil, err
	}
	log.Printf("Updated product count successfully for product  %v", product)
	return product, nil
}

// CreateProduct creates a new product
// func (s *ProductService) CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error) {
// 	return s.repo.CreateProduct(ctx, req)
// }

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}
