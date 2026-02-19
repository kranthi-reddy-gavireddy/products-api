package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	apperrors "products-api/app-errors"
	"products-api/dtos"
	"products-api/internal/models"
	"products-api/internal/repository"
	"products-api/logger"
	"time"

	"products-api/utils/cache"
	ctx "products-api/utils/context"

	"github.com/redis/go-redis/v9"
)

type IProductService interface {
	Create(context *ctx.Context, req *dtos.ProductRequest) (*dtos.ProductResponse, error)
	Update(ctx *ctx.Context, id string, req *dtos.ProductRequest) (*dtos.ProductResponse, error)
	UpdateCount(ctx context.Context, id string, sold int) (*models.Product, error)
	Get(ctx *ctx.Context, limit, offset int) ([]models.Product, error)
	Filter(ctx *ctx.Context, minPrice float64, maxPrice float64, category string, limit int, offset int, sortClause []dtos.SortClause) ([]models.Product, error)
	GetByID(ctx *ctx.Context, id string) (*dtos.ProductResponse, error)
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
		context.Logger.Errorf("Error creating product: %v", err)
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

	cacheData(context, res)
	context.Logger.Infof("Created product successfully: %+v", res)
	return res, nil
}

func (s *ProductService) Update(ctx *ctx.Context, id string, req *dtos.ProductRequest) (*dtos.ProductResponse, error) {

	product, err := s.repo.GetByID(ctx.Context(), id)
	if err != nil {
		ctx.Logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	product.Name = req.Name
	product.Price = req.Price
	product.Category = req.Category
	product.Quantity = req.Quantity
	product.SellerID = req.SellerID

	err = s.repo.Update(ctx.Context(), product)
	if err != nil {
		ctx.Logger.Errorf("Error updating product by ID %s: %v", id, err)
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

	cacheData(ctx, res)

	ctx.Logger.Infof("Updated product successfully: %+v", res)
	return res, nil
}

func (s *ProductService) UpdateCount(context context.Context, id string, sold int) (*models.Product, error) {

	logger := logger.New()

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

	ctx.Logger.Infof("Retrieving products with limit %d and offset %d", limit, offset)
	products, err := s.repo.GetAll(ctx.Context(), limit, offset)
	if err != nil {
		ctx.Logger.Errorf("Error retrieving products: %v", err)
		return nil, apperrors.DatabaseError()
	}

	ctx.Logger.Infof("Retrieved %d products %v", len(products), products)
	return products, nil
}

func (s *ProductService) Filter(context *ctx.Context, minPrice float64, maxPrice float64, category string, limit int, offset int, sortClause []dtos.SortClause) ([]models.Product, error) {

	var err error

	context.Logger.Infof("Filtering products with minPrice: %f, maxPrice: %f, category: %s, sorting %v limit: %d, offset: %d", minPrice, maxPrice, category, sortClause, limit, offset)
	products, err := s.repo.FilterProducts(context.Context(), minPrice, maxPrice, category, sortClause, limit, offset)
	if err != nil {
		context.Logger.Errorf("Error filtering products: %v", err)
		return nil, apperrors.DatabaseError()
	}

	context.Logger.Infof("Filtered products: %v", products)
	return products, nil
}

func (s *ProductService) GetByID(ctx *ctx.Context, id string) (*dtos.ProductResponse, error) {

	res := getCachedData(ctx, id)
	if res != nil {
		ctx.Logger.Infof("Returning cached product for ID %s: %v", id, res)
		return res, nil
	}

	product, err := s.repo.GetByID(ctx.Context(), id)
	if err != nil {
		ctx.Logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return nil, apperrors.DatabaseError()
	}

	res = &dtos.ProductResponse{
		ID: product.ID,
		ProductBase: dtos.ProductBase{
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
			Quantity: product.Quantity,
			SellerID: product.SellerID,
		},
	}

	cacheData(ctx, res)

	ctx.Logger.Infof("Retrieved product by ID %s: %v", id, product)
	return res, nil
}

func (s *ProductService) Delete(ctx *ctx.Context, id string) error {

	product, err := s.repo.GetByID(ctx.Context(), id)
	if err != nil {
		ctx.Logger.Errorf("Error retrieving product by ID %s: %v", id, err)
		return apperrors.DatabaseError()
	}

	ctx.Logger.Infof("Deleting product: %v", product)
	err = s.repo.DeleteProduct(ctx.Context(), id)
	if err != nil {
		ctx.Logger.Errorf("Error deleting product by ID %s: %v", id, err)
		return apperrors.DatabaseError()
	}

	removeCache(ctx, id)

	ctx.Logger.Infof("Deleted product by ID %s", id)
	return nil
}

func New(repo repository.IProductRepository) IProductService {
	return &ProductService{repo: repo}
}

func cacheData(ctx *ctx.Context, res *dtos.ProductResponse) {

	rdb := cache.New().Client

	bytes, err := json.Marshal(res)
	if err != nil {
		ctx.Logger.Errorf("Error marshaling product response for caching with ID %s: %v", res.ID, err)
		return
	}

	err = rdb.Set(ctx.Context(), res.ID, bytes, 24*time.Hour).Err()
	if err != nil {
		ctx.Logger.Errorf("Error caching product with ID %s: %v", res.ID, err)
	}
}

func removeCache(ctx *ctx.Context, id string) {

	rdb := cache.New().Client

	err := rdb.Del(ctx.Context(), id).Err()
	if err != nil {
		ctx.Logger.Errorf("Error removing cached product with ID %s: %v", id, err)
	}
}

func getCachedData(ctx *ctx.Context, id string) *dtos.ProductResponse {

	rdb := cache.New().Client

	var cachedRes dtos.ProductResponse
	bytes, err := rdb.Get(ctx.Context(), id).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			ctx.Logger.Infof("No cached product found with ID %s", id)
			return nil
		}

		ctx.Logger.Errorf("Error retrieving cached product with ID %s: %v", id, err)
		return nil
	}

	err = json.Unmarshal(bytes, &cachedRes)
	if err != nil {
		ctx.Logger.Errorf("Error unmarshaling cached product with ID %s: %v", id, err)
		return nil
	}

	ctx.Logger.Infof("Retrieved cached product with ID %s: %v", id, cachedRes)
	return &cachedRes
}
