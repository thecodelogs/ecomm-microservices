CREATE TYPE checkout_session_status AS ENUM (
    'pending',
    'processing',
    'confirmed',
    'cancelled',
    'refunded'
);

CREATE TABLE checkout_session (
  id                       UUID      PRIMARY KEY DEFAULT uuidv7(),
  user_id                  UUID          NOT NULL,
  status                   checkout_session_status NOT NULL DEFAULT 'pending',
  total_amount             NUMERIC(12,2) NOT NULL CHECK (total_amount >= 0),
  currency                 CHAR(3)       NOT NULL DEFAULT 'INR',
  payment_transaction_id   VARCHAR(255),
  metadata                 JSONB,
  created_at               TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at               TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_checkout_session_user_id  ON checkout_session (user_id);
CREATE INDEX idx_checkout_session_status   ON checkout_session (status);