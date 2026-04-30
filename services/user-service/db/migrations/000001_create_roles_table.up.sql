CREATE TABLE IF NOT EXISTS roles (
    id              UUID        PRIMARY KEY DEFAULT uuidv7(),
    name            VARCHAR(50)  NOT NULL UNIQUE,
    description     VARCHAR(255),
    is_default      BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Seed default roles
INSERT INTO roles (name, description, is_default) VALUES
    ('admin',    'Full access to all resources',              FALSE),
    ('customer', 'Default role for registered customers',     TRUE),
    ('seller',   'Can manage own products and view orders',   FALSE),
    ('support',  'Can view users and manage support tickets', FALSE);
