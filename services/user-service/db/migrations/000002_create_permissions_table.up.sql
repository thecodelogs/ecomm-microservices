CREATE TABLE IF NOT EXISTS permissions (
    id          SERIAL       PRIMARY KEY,
    resource    VARCHAR(100) NOT NULL,
    action      VARCHAR(50)  NOT NULL,
    description VARCHAR(255),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    UNIQUE (resource, action)
);

-- Seed common permissions
INSERT INTO permissions (resource, action, description) VALUES
    ('users',    'read',   'View user profiles'),
    ('users',    'write',  'Create and update users'),
    ('users',    'delete', 'Delete user accounts'),
    ('products', 'read',   'View products'),
    ('products', 'write',  'Create and update products'),
    ('products', 'delete', 'Delete products'),
    ('orders',   'read',   'View orders'),
    ('orders',   'write',  'Create and update orders'),
    ('orders',   'delete', 'Cancel or delete orders'),
    ('roles',    'read',   'View roles'),
    ('roles',    'write',  'Create and assign roles');
