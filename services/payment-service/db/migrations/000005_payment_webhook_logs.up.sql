CREATE TABLE payment_webhook_logs (
  id              UUID          PRIMARY KEY DEFAULT uuidv7(),
  gateway         VARCHAR(100)  NOT NULL,
  event_id        VARCHAR(255),                          -- gateway's own event ID (for dedup)
  event_type      VARCHAR(100)  NOT NULL,                -- e.g. 'payment.captured'
  payload         JSONB         NOT NULL,
  signature       VARCHAR(500),                          -- HMAC signature for verification
  is_verified     BOOLEAN       NOT NULL DEFAULT FALSE,  -- signature verified?
  processed       BOOLEAN       NOT NULL DEFAULT FALSE,
  processed_at    TIMESTAMPTZ,
  error           TEXT,
  received_at     TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
 
CREATE UNIQUE INDEX idx_webhook_event_dedup
  ON payment_webhook_logs (gateway, event_id)
  WHERE event_id IS NOT NULL;

CREATE INDEX idx_webhook_gateway    ON payment_webhook_logs (gateway);
CREATE INDEX idx_webhook_event_type ON payment_webhook_logs (event_type);
CREATE INDEX idx_webhook_processed  ON payment_webhook_logs (processed) WHERE processed = FALSE;
CREATE INDEX idx_webhook_received   ON payment_webhook_logs (received_at DESC);