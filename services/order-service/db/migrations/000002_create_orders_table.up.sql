ENUM "order_status" AS ENUM (
  'pending',
  'confirmed',
  'processing',
  'shipped',
  'out_for_delivery',
  'delivered',
  'cancelled',
  'returned',
  'refunded'
);

CREATE TABLE "order" (
  id                    UUID          PRIMARY KEY DEFAULT uuidv7(),
  display_order_id      VARCHAR(50)   NOT NULL UNIQUE,  -- e.g. ORD-20260430-00042
  checkout_session_id   UUID          NOT NULL REFERENCES checkout_session(id)
                                      ON DELETE RESTRICT,
  user_id               UUID          NOT NULL,
  seller_id             UUID,                           -- NULL = fulfilled by platform
  fulfillment_center_id UUID,
  status                order_status NOT NULL DEFAULT 'pending',
  subtotal              NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (subtotal >= 0),
  shipping_amount       NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (shipping_amount >= 0),
  tax_amount            NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (tax_amount >= 0),
  discount_amount       NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
  total_amount          NUMERIC(12,2) GENERATED ALWAYS AS (
                          subtotal + shipping_amount + tax_amount - discount_amount
                        ) STORED,
  notes                 TEXT,
  created_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_order_checkout_session_id   ON "order" (checkout_session_id);
CREATE INDEX idx_order_user_id               ON "order" (user_id);
CREATE INDEX idx_order_seller_id             ON "order" (seller_id);
CREATE INDEX idx_order_status                ON "order" (status);
CREATE INDEX idx_order_created_at            ON "order" (created_at DESC);