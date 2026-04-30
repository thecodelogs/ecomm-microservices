CREATE TYPE "order_item_status" AS ENUM (
  'active',
  'cancelled',
  'returned',
  'refunded'
);

CREATE TABLE order_item (
  id            UUID          PRIMARY KEY DEFAULT uuidv7(),
  order_id      UUID          NOT NULL REFERENCES "order"(id)
                              ON DELETE CASCADE,
  product_id    UUID          NOT NULL,
  variant_id    UUID,                                  -- size, colour, etc.
  product_name  VARCHAR(500)  NOT NULL,               -- snapshot
  sku           VARCHAR(100),
  quantity      INT           NOT NULL CHECK (quantity > 0),
  unit_price    NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
  total_price   NUMERIC(12,2) GENERATED ALWAYS AS (
                  quantity * unit_price
                ) STORED,
  status        order_item_status   NOT NULL DEFAULT 'active',
  metadata      JSONB,                                -- e.g. gift wrap, warranty
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_order_item_order_id    ON order_item (order_id);
CREATE INDEX idx_order_item_product_id  ON order_item (product_id);
CREATE INDEX idx_order_item_variant_id  ON order_item (variant_id);