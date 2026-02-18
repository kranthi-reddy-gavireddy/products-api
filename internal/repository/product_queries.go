package repository

var (
	CREATE_PRODUCT_QUERY = `INSERT INTO products (id, name, price, seller_id, quantity, category, created_at, updated_at)
	           VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, name, price, seller_id, quantity, category, created_at, updated_at`
	UPDATE_PRODUCT_QUERY = `UPDATE products SET name = $1, price = $2, seller_id = $3, quantity = $4, category = $5, updated_at = NOW() WHERE id = $6`
	COUNT_UPDATE_QUERY   = `UPDATE products SET quantity = $1, updated_at = NOW() WHERE id = $2`
	DELETE_PRODUCT_QUERY = `DELETE FROM products WHERE id = $1`
	GET_BY_ID_QUERY      = `SELECT id, name, price, seller_id, quantity, category, created_at, updated_at FROM products WHERE id = $1`
	GET_ALL_QUERY        = `SELECT id, name, price, seller_id, quantity, category, created_at, updated_at FROM products WHERE 1=1`
)
