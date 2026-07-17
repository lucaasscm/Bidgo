CREATE TABLE IF NOT EXISTS bids (
    id         UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID             NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    bidder_id  UUID             NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    bid_amount DOUBLE PRECISION NOT NULL CHECK (bid_amount > 0),
    created_at TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS bids_product_id_bid_amount_idx ON bids (product_id, bid_amount DESC);
CREATE INDEX IF NOT EXISTS bids_bidder_id_idx ON bids (bidder_id);

---- create above / drop below ----

DROP INDEX IF EXISTS bids_bidder_id_idx;
DROP INDEX IF EXISTS bids_product_id_bid_amount_idx;
DROP TABLE IF EXISTS bids;
