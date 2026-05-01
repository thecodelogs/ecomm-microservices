CREATE TYPE payment_type AS ENUM (
    'card', 'upi', 'netbanking', 'wallet', 'cod', 'bnpl'
);

CREATE TABLE payment_method (
  id               UUID          PRIMARY KEY DEFAULT uuidv7(),
  user_ref         UUID          NOT NULL,
  type             payment_type  NOT NULL,
  provider         VARCHAR(100)  NOT NULL,  
  provider_token   VARCHAR(500)  NOT NULL,              
  last4            VARCHAR(4),                          
  brand            VARCHAR(50),                         
  bank_name        VARCHAR(100),                        
  upi_id           VARCHAR(255),                        
  is_default       BOOLEAN       NOT NULL DEFAULT FALSE,
  is_active        BOOLEAN       NOT NULL DEFAULT TRUE,
  expires_at       TIMESTAMPTZ,                         
  metadata         JSONB,
  created_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);
 
CREATE INDEX idx_pm_user_ref    ON payment_method (user_ref);
CREATE INDEX idx_pm_type        ON payment_method (type);
CREATE INDEX idx_pm_is_default  ON payment_method (user_ref, is_default) WHERE is_default = TRUE;