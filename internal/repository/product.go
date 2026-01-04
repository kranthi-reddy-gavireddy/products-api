package repository

import (
	"context"
	"database/sql"
	"products-api/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := "SELECT id, name, price, seller_id, quantity, created_at, updated_at FROM products"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.SellerID, &p.Quantity, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Create(ctx context.Context, req models.Product) (*models.Product, error) {
	query := `INSERT INTO products (id, name, price, seller_id, quantity, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, name, price, seller_id, quantity, created_at, updated_at`
	var p models.Product
	req.SetID()
	err := r.db.QueryRowContext(ctx, query, req.ID, req.Name, req.Price, req.SellerID, req.Quantity).Scan(
		&p.ID, &p.Name, &p.Price, &p.SellerID, &p.Quantity, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) UpdateProductCount(ctx context.Context, product *models.Product, sold int) (*models.Product, error) {
	product.Quantity -= sold
	query := `UPDATE products SET quantity = $1, updated_at = NOW() WHERE id = $2`
	r.db.QueryRowContext(ctx, query, product.Quantity, product.ID)
	return product, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id string) error {
	query := "DELETE FROM products WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *ProductRepository) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	query := "SELECT id, name, description, price, quantity, created_at, updated_at FROM products WHERE id = $1"
	var p models.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Quantity)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
