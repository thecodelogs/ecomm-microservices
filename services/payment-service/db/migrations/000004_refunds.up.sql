CREATE TYPE refund_status AS ENUM (
  'pending',
  'processing',
  'succeeded',
  'failed',
  'cancelled'
);

CREATE TABLE refunds (
  id                  UUID          PRIMARY KEY DEFAULT uuidv7(),
  payment_id          UUID          NOT NULL REFERENCES payments(id)
                                    ON DELETE RESTRICT,
  order_item_ref      UUID,

  amount              NUMERIC(12,2) NOT NULL CHECK (amount > 0),
  currency            CHAR(3)       NOT NULL DEFAULT 'INR',
  reason              TEXT,
  status              refund_status NOT NULL DEFAULT 'pending',
  gateway_refund_id   VARCHAR(255)  UNIQUE,
  gateway_response    JSONB,
  initiated_by        VARCHAR(100)  NOT NULL,  -- 'user', 'admin:<id>', 'system'
  refunded_at         TIMESTAMPTZ,
  created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refund_payment    ON refunds (payment_id);
CREATE INDEX idx_refund_status     ON refunds (status);
CREATE INDEX idx_refund_gateway_id ON refunds (gateway_refund_id);