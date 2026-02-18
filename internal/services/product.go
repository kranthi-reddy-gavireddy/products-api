package services

import (
	"context"
	"fmt"
	apperrors "products-api/app-errors"
	"products-api/dtos"
	"products-api/internal/models"
	"products-api/internal/repository"
	"products-api/logger"

	ctx "products-api/utils/context"
)

type IProductService interface {
	Create(context *ctx.Context, req *dtos.ProductRequest) (*dtos.ProductResponse, error)
	Update(ctx *ctx.Context, id string, req *dtos.ProductRequest) (*dtos.ProductResponse, error)
	UpdateCount(ctx context.Context, id string, sold int) (*models.Product, error)
	Get(ctx *ctx.Context, limit, offset int) ([]models.Product, error)
	Filter(ctx *ctx.Context, minPrice float64, maxPrice float64, category string, limit int, offset int, sortClause []dtos.SortClause) ([]models.Product, error)
	GetByID(ctx *ctx.Context, id string) (*models.Product, error)
	Delete(ctx *ctx.Context, id string) error
}

// ProductService handles product business logic
type ProductService struct {
	repo repository.IProductRepository
}

func (s *ProductService) Create(context *ctx.Context, req *dtos.ProductRequest) (*dtos.ProductResponse, error) {

	var err error

	product := &models.Product{
		Name:     req.Name,
		Price:    req.Price,
		Category: req.Category,
		Quantity: req.Quantity,
		SellerID: req.SellerID,
	}
	product.SetID()

	err = s.repo.Create(context.Context(), product)
	if err != nil {
		logger.Errorf("Error creating product: %v", err)
		return nil, apperrors.DatabaseError()
	}

	res := &dtos.ProductResponse{
		ID: product.ID,
		ProductBase: dtos.ProductBase{
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

func (s *ProductService) Update(ctx *ctx.Context, id string, req *dtos.ProductRequest) (*dtos.ProductResponse, error) {

	product, err := s.repo.GetByID(ctx.Context(), id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Category = req.Category
	product.Quantity = req.Quantity
	product.SellerID = req.SellerID

	err = s.repo.Update(ctx.Context(), product)
	if err != nil {
		logger.Errorf("Error updating product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	res := &dtos.ProductResponse{
		ID: product.ID,
		ProductBase: dtos.ProductBase{
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
			Quantity: product.Quantity,
			SellerID: product.SellerID,
		},
	}

	logger.Infof("Updated product successfully: %+v", res)
	return res, nil
}

func (s *ProductService) UpdateCount(context context.Context, id string, sold int) (*models.Product, error) {

	product, err := s.repo.GetByID(context, id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	if product.Quantity < sold {
		return nil, apperrors.BadRequestError(fmt.Sprintf(apperrors.INSUFFICIENT_QUANTITY, id))
	}

	logger.Infof("Updating product count for product ID %s, sold: %d", id, sold)
	err = s.repo.UpdateProductCount(context, product, sold)
	if err != nil {
		logger.Errorf("Error updating product count for product ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Updated product count successfully for product  %v", product)
	return product, nil
}

// GetProducts retrieves all products
func (s *ProductService) Get(ctx *ctx.Context, limit, offset int) ([]models.Product, error) {

	var err error

	logger.Infof("Retrieving products with limit %d and offset %d", limit, offset)
	products, err := s.repo.GetAll(ctx.Context(), limit, offset)
	if err != nil {
		logger.Errorf("Error retrieving products: %v", err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Retrieved %d products %v", len(products), products)
	return products, nil
}

func (s *ProductService) Filter(context *ctx.Context, minPrice float64, maxPrice float64, category string, limit int, offset int, sortClause []dtos.SortClause) ([]models.Product, error) {

	var err error

	logger.Infof("Filtering products with minPrice: %f, maxPrice: %f, category: %s, sorting %v limit: %d, offset: %d", minPrice, maxPrice, category, sortClause, limit, offset)
	products, err := s.repo.FilterProducts(context.Context(), minPrice, maxPrice, category, sortClause, limit, offset)
	if err != nil {
		logger.Errorf("Error filtering products: %v", err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Filtered products: %v", products)
	return products, nil
}

func (s *ProductService) GetByID(ctx *ctx.Context, id string) (*models.Product, error) {

	product, err := s.repo.GetByID(ctx.Context(), id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	logger.Infof("Retrieved product by ID %s: %v", id, product)
	return product, nil
}

func (s *ProductService) Delete(ctx *ctx.Context, id string) error {

	product, err := s.repo.GetByID(ctx.Context(), id)
	if err != nil {
		logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return apperrors.DatabaseError()
	}

	logger.Infof("Deleting product: %v", product)
	err = s.repo.DeleteProduct(ctx.Context(), id)
	if err != nil {
		logger.Errorf("Error deleting product by ID %s: %v", id, err)
		return apperrors.DatabaseError()
	}

	logger.Infof("Deleted product by ID %s", id)
	return nil
}

func New(repo repository.IProductRepository) IProductService {
	return &ProductService{repo: repo}
}
