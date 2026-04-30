CREATE TABLE IF NOT EXISTS addresses (
    id            UUID         PRIMARY KEY DEFAULT uuidv7(),
    user_id       UUID         NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    label         VARCHAR(50)  NOT NULL DEFAULT 'home',
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city          VARCHAR(100) NOT NULL,
    state         VARCHAR(100) NOT NULL,
    postal_code   VARCHAR(20)  NOT NULL,
    country       VARCHAR(100) NOT NULL DEFAULT 'India',
    landmark      VARCHAR(255),
    phone         VARCHAR(20),
    is_default    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_addresses_user_id ON addresses (user_id);

-- Ensure only one default address per user
CREATE UNIQUE INDEX idx_addresses_user_default
    ON addresses (user_id) WHERE is_default = TRUE;
