-- Insert default roles
INSERT INTO roles (id, name, description)
VALUES 
    ('832b84eb-bdc9-467f-9457-3a830af96d10', 'admin', 'Administrator role'),
    ('e138a377-6bb9-4b6e-a7f4-d57be2c9e782', 'customer', 'Standard customer role')
ON CONFLICT (name) DO NOTHING;

-- Insert default admin user (password: admin123)
INSERT INTO users (
    id, email, password_hash, first_name, last_name, status, is_email_verified
) VALUES (
    '5d8a8b13-9a3b-4b2a-89a5-7b56d8c0b91e',
    'admin@ecomm.com',
    '$2a$10$6x4dM2sjA0XYYzUkjo8reuDzQ2Muzu.K6pBQbvjHCp.oSnIsGxZDy',
    'System',
    'Admin',
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- Assign admin role to the admin user
INSERT INTO user_roles (user_id, role_id)
VALUES (
    '5d8a8b13-9a3b-4b2a-89a5-7b56d8c0b91e',
    '832b84eb-bdc9-467f-9457-3a830af96d10'
) ON CONFLICT DO NOTHING;
