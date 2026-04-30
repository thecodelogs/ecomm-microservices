CREATE TABLE IF NOT EXISTS brands
(
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    name        VARCHAR(255) NOT NULL,
    description text,
    image_url   VARCHAR(255),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_brands_name ON brands(name);