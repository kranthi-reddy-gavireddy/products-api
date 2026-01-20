package services

import (
	"context"
	"errors"
	"log"
	"products-api/internal/models"
	"products-api/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	repo *repository.ProductRepository
}

func (s *ProductService) FilterProducts(ctx context.Context, minPrice float64, maxPrice float64, category string, limit int, offset int) ([]models.Product, error) {
	var err error
	log.Printf("Filtering products with minPrice: %f, maxPrice: %f, category: %s, limit: %d, offset: %d", minPrice, maxPrice, category, limit, offset)
	products, err := s.repo.FilterProducts(ctx, minPrice, maxPrice, category, limit, offset)
	if err != nil {
		log.Printf("Error filtering products: %v", err)
		return nil, err
	}
	log.Printf("Filtered products: %v", products)
	return products, nil
}

func (s *ProductService) Create(context context.Context, product *models.Product) error {
	var err error
	err = s.repo.Create(context, product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
	}
	return err
}

// GetProducts retrieves all products
func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	var err error
	log.Printf("Retrieving products with limit %d and offset %d", limit, offset)
	products, err := s.repo.GetAll(ctx, limit, offset)

	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		return nil, err
	}
	log.Printf("Retrieved %d products %v", len(products), products)
	return products, nil
}

func (s *ProductService) UpdateProductCount(ctx context.Context, id string, sold int) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieving product by ID: %v", err)
		return nil, err
	}
	if product.Quantity < sold {
		log.Printf("Insufficient product quantity: available %d, requested %d", product.Quantity, sold)
		return nil, errors.New("Insufficient product quantity to fulfill the order")
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

func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieving product by ID: %v", err)
		return nil, err
	}
	log.Printf("Retrieved product by ID %s: %v", id, product)
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Error retrieving product by ID before deletion: %v", err)
		return err
	}
	log.Printf("Deleting product: %v", product)
	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		log.Printf("Error deleting product by ID: %v", err)
		return err
	}
	log.Printf("Deleted product by ID %s", id)
	return nil
}

// CreateProduct creates a new product
// func (s *ProductService) CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error) {
// 	return s.repo.CreateProduct(ctx, req)
// }

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}
