CREATE TABLE IF NOT EXISTS product_categories
(
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    product_id  UUID        NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id UUID        NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);