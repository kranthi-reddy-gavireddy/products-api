package services

import (
	"context"
	"fmt"
	apperrors "products-api/app-errors"
	"products-api/dtos"
	"products-api/internal/models"
	"products-api/internal/repository"
	"products-api/logger"
)

// ProductService handles product business logic
type ProductService struct {
	repo *repository.ProductRepository
}

func (s *ProductService) FilterProducts(ctx context.Context, minPrice float64, maxPrice float64, category string, limit int, offset int, sortClause []dtos.SortClause) ([]models.Product, error) {

	var err error

	logger.Infof("Filtering products with minPrice: %f, maxPrice: %f, category: %s, sorting %v limit: %d, offset: %d", minPrice, maxPrice, category, sortClause, limit, offset)
	products, err := s.repo.FilterProducts(ctx, minPrice, maxPrice, category, sortClause, limit, offset)
	if err != nil {
		logger.Errorf("Error filtering products: %v", err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Filtered products: %v", products)
	return products, nil
}

func (s *ProductService) Create(context context.Context, req *dtos.ProductRequest) (*dtos.ProductResponse, error) {

	var err error

	product := &models.Product{
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
		Quantity: req.Quantity,
	}
	product.SetID()

	err = s.repo.Create(context, product)
	if err != nil {
		logger.Errorf("Error creating product: %v", err)
		return nil, apperrors.DatabaseError()
	}

	res := &dtos.ProductResponse{
		ID: product.ID,
		ProductRequest: dtos.ProductRequest{
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
			Quantity: product.Quantity,
			SellerID: product.SellerID,
		},
	}

	logger.Infof("Created product successfully: %+v", res)
	return res, nil
}

// GetProducts retrieves all products
func (s *ProductService) GetProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {

	var err error

	logger.Infof("Retrieving products with limit %d and offset %d", limit, offset)
	products, err := s.repo.GetAll(ctx, limit, offset)
	if err != nil {
		logger.Errorf("Error retrieving products: %v", err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Retrieved %d products %v", len(products), products)
	return products, nil
}

func (s *ProductService) UpdateProductCount(ctx context.Context, id string, sold int) (*models.Product, error) {

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	if product.Quantity < sold {
		return nil, apperrors.BadRequestError(fmt.Sprintf(apperrors.INSUFFICIENT_QUANTITY, id))
	}

	logger.Infof("Updating product count for product ID %s, sold: %d", id, sold)
	err = s.repo.UpdateProductCount(ctx, product, sold)
	if err != nil {
		logger.Errorf("Error updating product count for product ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Updated product count successfully for product  %v", product)
	return product, nil
}

func (s *ProductService) GetByID(ctx context.Context, id string) (*models.Product, error) {

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Retrieved product by ID %s: %v", id, product)
	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id string) error {

	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return apperrors.DatabaseError()
	}

	logger.Infof("Deleting product: %v", product)
	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		logger.Errorf("Error deleting product by ID %s: %v", id, err)
		return apperrors.DatabaseError()
	}

	logger.Infof("Deleted product by ID %s", id)
	return nil
}

// CreateProduct creates a new product
// func (s *ProductService) CreateProduct(ctx context.Context, req models.ProductCreateRequest) (*models.Product, error) {
// 	return s.repo.CreateProduct(ctx, req)
// }

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}
