CREATE TABLE IF NOT EXISTS products (
    id           UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    product_name VARCHAR(255)     NOT NULL,
    base_price   DOUBLE PRECISION NOT NULL CHECK (base_price >= 0),
    seller_id    UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    auction_end  TIMESTAMPTZ      NOT NULL,
    is_sold      BOOLEAN          NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS products_seller_id_idx ON products (seller_id);

---- create above / drop below ----

DROP INDEX IF EXISTS products_seller_id_idx;
DROP TABLE IF EXISTS products;
