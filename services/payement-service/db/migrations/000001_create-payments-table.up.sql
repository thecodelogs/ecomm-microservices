CREATE TYPE payment_status AS ENUM (
    'pending',
    'authorized',
    'captured',
    'failed',
    'refunded',
    'partially_refunded'
);

CREATE TABLE payment (
  id                  UUID          PRIMARY KEY DEFAULT uuidv7(),
  checkout_session_id UUID          NOT NULL REFERENCES checkout_session(id)
                                    ON DELETE RESTRICT,
  gateway             VARCHAR(100)  NOT NULL,          -- e.g. 'razorpay', 'stripe'
  transaction_id      VARCHAR(255)  UNIQUE,
  status              payment_status   NOT NULL DEFAULT 'pending',
  amount              NUMERIC(12,2) NOT NULL CHECK (amount >= 0),
  currency            CHAR(3)       NOT NULL DEFAULT 'INR',
  gateway_response    JSONB,                           -- raw gateway payload
  paid_at             TIMESTAMPTZ,
  created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_payment_checkout_session_id ON payment (checkout_session_id);
CREATE INDEX idx_payment_transaction_id      ON payment (transaction_id);
CREATE INDEX idx_payment_status              ON payment (status);