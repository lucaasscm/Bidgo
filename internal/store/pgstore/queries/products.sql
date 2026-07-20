-- name: CreateProduct :one
INSERT INTO products (product_name, base_price, seller_id, auction_end)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: GetProductById :one
SELECT id, product_name, base_price, seller_id, auction_end, is_sold, created_at, updated_at
FROM products
WHERE id = $1;

-- name: GetProductByIdForUpdate :one
SELECT id, product_name, base_price, seller_id, auction_end, is_sold, created_at, updated_at
FROM products
WHERE id = $1
FOR UPDATE;

-- name: ListProducts :many
SELECT id, product_name, base_price, seller_id, auction_end, is_sold, created_at, updated_at
FROM products
ORDER BY created_at DESC;

-- name: ListProductsBySeller :many
SELECT id, product_name, base_price, seller_id, auction_end, is_sold, created_at, updated_at
FROM products
WHERE seller_id = $1
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET product_name = $2,
    base_price   = $3,
    auction_end  = $4,
    is_sold      = $5,
    updated_at   = NOW()
WHERE id = $1
RETURNING id, product_name, base_price, seller_id, auction_end, is_sold, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
