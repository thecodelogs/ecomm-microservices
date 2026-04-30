CREATE TABLE IF NOT EXISTS categories (
    id          UUID        PRIMARY KEY DEFAULT uuidv7(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id   UUID NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);