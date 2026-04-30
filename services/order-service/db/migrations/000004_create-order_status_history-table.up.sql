CREATE TABLE order_status_history (
  id           UUID         PRIMARY KEY DEFAULT uuidv7(),
  order_id     UUID         NOT NULL REFERENCES "order"(id)
                            ON DELETE CASCADE,
  from_status  VARCHAR(50),                          -- NULL on first insert
  to_status    VARCHAR(50)  NOT NULL,
  reason       TEXT,                                 -- e.g. 'payment failed'
  changed_by   VARCHAR(100) NOT NULL,                -- 'user', 'system', 'admin:<id>'
  metadata     JSONB,
  changed_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_order_status_history_order_id    ON order_status_history (order_id);
CREATE INDEX idx_order_status_history_changed_at  ON order_status_history (changed_at DESC);