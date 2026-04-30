CREATE TABLE IF NOT EXISTS products
(
    id            UUID        PRIMARY KEY DEFAULT uuidv7(),
    name          VARCHAR(255) NOT NULL,
    description   text,
    price         numeric(10, 2) NOT NULL,
    brand_id      UUID,
    stock         integer DEFAULT 0,
    image_url     text,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE SET NULL
);

CREATE INDEX idx_products_name ON products(name);
