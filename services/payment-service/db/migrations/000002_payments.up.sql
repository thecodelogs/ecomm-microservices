DO $$ BEGIN
  CREATE TYPE payment_status AS ENUM (
    'pending', 'authorized', 'captured', 'failed',
    'cancelled', 'refunded', 'partially_refunded'
  );
EXCEPTION WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE payments (
  id                    UUID          PRIMARY KEY DEFAULT uuidv7(),
 
  checkout_session_ref  UUID          NOT NULL,
  user_ref              UUID          NOT NULL,
 
  payment_method_id     UUID          REFERENCES payment_method(id)
                                      ON DELETE SET NULL,
 
  gateway               VARCHAR(100)  NOT NULL,
  gateway_order_id      VARCHAR(255),
  transaction_id        VARCHAR(255)  UNIQUE,
  status                payment_status  NOT NULL DEFAULT 'pending',
  amount                NUMERIC(12,2) NOT NULL CHECK (amount > 0),
  currency              CHAR(3)       NOT NULL DEFAULT 'INR',
  amount_refunded       NUMERIC(12,2) NOT NULL DEFAULT 0
                                      CHECK (amount_refunded >= 0),
 
  gateway_response      JSONB,
 
  failure_code          VARCHAR(100),
  failure_message       TEXT,
 
  paid_at               TIMESTAMPTZ,
  created_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
 
  CONSTRAINT amount_refunded_lte_amount
    CHECK (amount_refunded <= amount)
);
 
CREATE INDEX idx_payment_session_ref   ON payments (checkout_session_ref);
CREATE INDEX idx_payment_user_ref      ON payments (user_ref);
CREATE INDEX idx_payment_transaction   ON payments (transaction_id);
CREATE INDEX idx_payment_status        ON payments (status);
CREATE INDEX idx_payment_gateway       ON payments (gateway);
CREATE INDEX idx_payment_created_at    ON payments (created_at DESC);
 