CREATE TYPE payment_attempt_status AS ENUM (
  'initiated',
  'pending',
  'authorized',
  'captured',
  'failed',
  'cancelled',
  'timeout'
);

CREATE TABLE payment_attempts (
  id                UUID          PRIMARY KEY DEFAULT uuidv7(),
  payment_id        UUID          NOT NULL REFERENCES payments(id)
                                  ON DELETE CASCADE,
  attempt_number    SMALLINT      NOT NULL DEFAULT 1 CHECK (attempt_number > 0),
  gateway           VARCHAR(100)  NOT NULL,
  status            payment_attempt_status NOT NULL DEFAULT 'pending',
  amount            NUMERIC(12,2) NOT NULL CHECK (amount > 0),
  gateway_response  JSONB,
  failure_code      VARCHAR(100),
  failure_message   TEXT,
  ip_address        INET,
  user_agent        TEXT,
  attempted_at      TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
 
  UNIQUE (payment_id, attempt_number)
);
 


