package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// RunMigrations creates the necessary database tables
func RunMigrations(db *sql.DB) error {
	log.Println("Running database migrations...")

	// Create products table
	productsTable := `
	CREATE TABLE IF NOT EXISTS products (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		price DECIMAL(10,2) NOT NULL,
		seller_id VARCHAR(255),
		quantity INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE
	);`

	// Create indexes
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);`,
		`CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);`,
	}

	ctx := context.Background()

	// Execute table creation
	if _, err := db.ExecContext(ctx, productsTable); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	// Execute index creation
	for _, index := range indexes {
		if _, err := db.ExecContext(ctx, index); err != nil {
			log.Printf("Warning: failed to create index: %v", err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
