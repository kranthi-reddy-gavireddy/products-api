package repository

import (
	"context"
	"database/sql"

	"github.com/kranthi-reddy-gavireddy/products-api.git/dtos"
	"github.com/kranthi-reddy-gavireddy/products-api.git/internal/models"
)

type IProductRepository interface {
	Create(ctx context.Context, req *models.Product) error
	Update(ctx context.Context, req *models.Product) error
	UpdateProductCount(ctx context.Context, product *models.Product, sold int) error
	GetAll(ctx context.Context, limit, offset int) ([]models.Product, error)
	GetByID(ctx context.Context, id string) (*models.Product, error)
	FilterProducts(ctx context.Context, minPrice, maxPrice float64, category string, sortClauses []dtos.SortClause, limit, offset int) (*models.ListProducts, error)
	DeleteProduct(ctx context.Context, id string) error
}

type ProductRepository struct {
	db *sql.DB
}

func (r *ProductRepository) Create(ctx context.Context, req *models.Product) error {
	// query := `INSERT INTO products (id, name, price, seller_id, quantity, category, created_at, updated_at)
	//           VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, name, price, seller_id, quantity, category, created_at, updated_at`
	var p models.Product

	err := r.db.QueryRowContext(ctx, CREATE_PRODUCT_QUERY, req.ID, req.Name, req.Price, req.SellerID, req.Quantity, req.Category).Scan(
		&p.ID, &p.Name, &p.Price, &p.SellerID, &p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Update(ctx context.Context, req *models.Product) error {

	_, err := r.db.ExecContext(ctx, UPDATE_PRODUCT_QUERY, req.Name, req.Price, req.SellerID, req.Quantity, req.Category, req.ID)
	return err
}

func (r *ProductRepository) UpdateProductCount(ctx context.Context, product *models.Product, sold int) error {

	product.Quantity -= sold

	_, err := r.db.ExecContext(ctx, COUNT_UPDATE_QUERY, product.Quantity, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) GetAll(ctx context.Context, limit, offset int) ([]models.Product, error) {

	query := "SELECT id, name, price, seller_id, quantity, category, created_at, updated_at FROM products ORDER BY created_at DESC LIMIT $1 OFFSET $2"

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.SellerID, &p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetByID(ctx context.Context, id string) (*models.Product, error) {
	//query := "SELECT id, name, price, seller_id, quantity, category, created_at, updated_at FROM products WHERE id = $1"

	var p models.Product

	err := r.db.QueryRowContext(ctx, GET_BY_ID_QUERY, id).Scan(&p.ID, &p.Name, &p.Price, &p.SellerID, &p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProductRepository) FilterProducts(ctx context.Context, minPrice, maxPrice float64, category string, sortClauses []dtos.SortClause, limit, offset int) (*models.ListProducts, error) {

	filterparameters := " AND price >= $1 AND price <= $2"
	if category != "" {
		filterparameters += " AND category = $5"
	}

	sortOrder := sortQueryGenerator(sortClauses)

	query := GET_ALL_QUERY + filterparameters + sortOrder + " LIMIT $3 OFFSET $4"
	totalCountQuery := GET_PRODUCTS_COUNT_QUERY + filterparameters

	var (
		rows       *sql.Rows
		err        error
		totalCount int
	)

	if category != "" {
		rows, err = r.db.QueryContext(ctx, query, minPrice, maxPrice, limit, offset, category)
		r.db.QueryRowContext(ctx, totalCountQuery, minPrice, maxPrice, category).Scan(&totalCount)
	} else {
		rows, err = r.db.QueryContext(ctx, query, minPrice, maxPrice, limit, offset)
		r.db.QueryRowContext(ctx, totalCountQuery, minPrice, maxPrice).Scan(&totalCount)
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product

		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.SellerID, &p.Quantity, &p.Category, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	productsList := &models.ListProducts{
		Products:   products,
		TotalCount: totalCount,
	}

	return productsList, rows.Err()
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {

	_, err := r.db.ExecContext(ctx, DELETE_PRODUCT_QUERY, id)
	return err
}

func sortQueryGenerator(sortClauses []dtos.SortClause) string {

	if len(sortClauses) == 0 {
		return " ORDER BY created_at DESC "
	}

	query := " ORDER BY "
	for idx, clause := range sortClauses {
		if idx > 0 {
			query += " , "
		}
		query += clause.Field + " " + clause.Direction
	}

	return query
}

func New(db *sql.DB) IProductRepository {
	return &ProductRepository{db: db}
}
