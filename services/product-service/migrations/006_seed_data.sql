INSERT INTO brands (id, name, description, is_active) VALUES ('b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'BrandX', 'Default brand for seed data', true) ON CONFLICT (id) DO NOTHING;
INSERT INTO categories (id, slug, name, description, is_active) VALUES ('d4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'electronics', 'Electronics', 'Electronic devices and accessories', true) ON CONFLICT (slug) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('be5a8e42-aee7-476b-8b04-69f4e94f7cbd', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-1', 'Awesome Product 1', 'This is the detailed description for Awesome Product 1.', 'Short description 1', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('3d730d6d-af42-4290-860f-dc352afc0e21', 'be5a8e42-aee7-476b-8b04-69f4e94f7cbd', 'SKU-1', 'Default Variant', 11.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('3d730d6d-af42-4290-860f-dc352afc0e21', 101, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('a0632be3-f363-4841-b80a-2660b73fee2f', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-2', 'Awesome Product 2', 'This is the detailed description for Awesome Product 2.', 'Short description 2', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('eb7f2796-62e9-481c-8232-e2875da4884c', 'a0632be3-f363-4841-b80a-2660b73fee2f', 'SKU-2', 'Default Variant', 12.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('eb7f2796-62e9-481c-8232-e2875da4884c', 102, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('1ac737d7-9229-49e7-8e4b-a6cd939a8acd', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-3', 'Awesome Product 3', 'This is the detailed description for Awesome Product 3.', 'Short description 3', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('14aae7eb-3094-44a1-8770-e56e23811f12', '1ac737d7-9229-49e7-8e4b-a6cd939a8acd', 'SKU-3', 'Default Variant', 13.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('14aae7eb-3094-44a1-8770-e56e23811f12', 103, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('d19105ae-4b16-4823-814d-853cf4f7f013', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-4', 'Awesome Product 4', 'This is the detailed description for Awesome Product 4.', 'Short description 4', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('6846b022-2595-4d91-95cc-90d016b4b461', 'd19105ae-4b16-4823-814d-853cf4f7f013', 'SKU-4', 'Default Variant', 14.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('6846b022-2595-4d91-95cc-90d016b4b461', 104, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('48ad2098-cf96-4dd1-99ec-99f5e6a7229d', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-5', 'Awesome Product 5', 'This is the detailed description for Awesome Product 5.', 'Short description 5', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('0b638139-0222-41cd-ac86-9f3244b23c92', '48ad2098-cf96-4dd1-99ec-99f5e6a7229d', 'SKU-5', 'Default Variant', 15.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('0b638139-0222-41cd-ac86-9f3244b23c92', 105, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('f24f835f-927d-4636-9ea6-fd04adfffe66', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-6', 'Awesome Product 6', 'This is the detailed description for Awesome Product 6.', 'Short description 6', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('2246aa9b-fdf7-4910-bb7e-30e52a0d9fbe', 'f24f835f-927d-4636-9ea6-fd04adfffe66', 'SKU-6', 'Default Variant', 16.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('2246aa9b-fdf7-4910-bb7e-30e52a0d9fbe', 106, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('84600591-ab4f-4dd2-b2b5-4d4cc6c5da65', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-7', 'Awesome Product 7', 'This is the detailed description for Awesome Product 7.', 'Short description 7', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('7c475a7b-c950-4269-8735-225411b74f44', '84600591-ab4f-4dd2-b2b5-4d4cc6c5da65', 'SKU-7', 'Default Variant', 17.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('7c475a7b-c950-4269-8735-225411b74f44', 107, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('d63e2cbb-9237-4c4a-8b27-f26af7f89140', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-8', 'Awesome Product 8', 'This is the detailed description for Awesome Product 8.', 'Short description 8', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('c83cf447-6777-497a-9202-f24f16af3869', 'd63e2cbb-9237-4c4a-8b27-f26af7f89140', 'SKU-8', 'Default Variant', 18.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('c83cf447-6777-497a-9202-f24f16af3869', 108, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('9323e32b-7f59-4ff7-bee0-6261aee8925c', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-9', 'Awesome Product 9', 'This is the detailed description for Awesome Product 9.', 'Short description 9', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('7342819a-037d-41df-b4a4-9425cbfc5b62', '9323e32b-7f59-4ff7-bee0-6261aee8925c', 'SKU-9', 'Default Variant', 19.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('7342819a-037d-41df-b4a4-9425cbfc5b62', 109, 0) ON CONFLICT (variant_id) DO NOTHING;

INSERT INTO products (id, category_id, slug, name, description, short_description, brand, brand_id, status) VALUES ('a28035f1-84a4-4321-9e06-42df7860b209', 'd4a9078e-814a-4a8e-a716-1bb29aa40e6c', 'product-10', 'Awesome Product 10', 'This is the detailed description for Awesome Product 10.', 'Short description 10', 'BrandX', 'b77a7f9a-1111-4b8a-9b1b-1234567890ab', 'active') ON CONFLICT (slug) DO NOTHING;
INSERT INTO variants (id, product_id, sku, name, price, is_active) VALUES ('2d1f1fb1-5af2-4956-8ecc-232b506afba0', 'a28035f1-84a4-4321-9e06-42df7860b209', 'SKU-10', 'Default Variant', 20.00, true) ON CONFLICT (sku) DO NOTHING;
INSERT INTO inventory (variant_id, quantity_on_hand, quantity_reserved) VALUES ('2d1f1fb1-5af2-4956-8ecc-232b506afba0', 110, 0) ON CONFLICT (variant_id) DO NOTHING;
