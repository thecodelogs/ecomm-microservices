CREATE TABLE IF NOT EXISTS role_permissions (
    role_id       UUID NOT NULL REFERENCES roles(id)       ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (role_id, permission_id)
);

-- Admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, p.id
    FROM roles r, permissions p
    WHERE r.name = 'admin';

-- Customer: read products, read/write own orders
INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, p.id
    FROM roles r, permissions p
    WHERE r.name = 'customer'
      AND (
          (p.resource = 'products' AND p.action = 'read')
       OR (p.resource = 'orders'   AND p.action IN ('read', 'write'))
       OR (p.resource = 'users'    AND p.action = 'read')
      );

-- Seller: manage own products, view orders
INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, p.id
    FROM roles r, permissions p
    WHERE r.name = 'seller'
      AND (
          (p.resource = 'products' AND p.action IN ('read', 'write', 'delete'))
       OR (p.resource = 'orders'   AND p.action = 'read')
       OR (p.resource = 'users'    AND p.action = 'read')
      );

-- Support: read users & orders
INSERT INTO role_permissions (role_id, permission_id)
    SELECT r.id, p.id
    FROM roles r, permissions p
    WHERE r.name = 'support'
      AND (
          (p.resource = 'users'  AND p.action = 'read')
       OR (p.resource = 'orders' AND p.action IN ('read', 'write'))
      );
