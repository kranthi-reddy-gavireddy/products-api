#!/bin/bash
# Wait for PostgreSQL to be ready
while ! pg_isready -h psql_bp -p 5432 -U test; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "PostgreSQL is ready. Creating tables..."

# Create products table
psql -h psql_bp -U test -d test -c "
CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    seller_id VARCHAR(255),
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create index on seller_id for better query performance
CREATE INDEX IF NOT EXISTS idx_products_seller_id ON products(seller_id);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);
"

echo "Database tables created successfully"