CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_created_at ON products(created_at DESC);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_products_updated_at ON products(updated_at DESC);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

INSERT INTO products (name, price, description) VALUES
    ('Laptop Dell XPS 13', 1299.99, 'High-performance ultrabook with 11th Gen Intel Core processor'),
    ('Apple MacBook Pro 14"', 1999.00, 'M2 Pro chip with 16GB RAM and 512GB SSD'),
    ('Wireless Mouse Logitech MX Master 3', 99.99, 'Advanced wireless mouse with ergonomic design'),
    ('Mechanical Keyboard Keychron K2', 79.99, 'Compact 75% wireless mechanical keyboard'),
    ('Monitor LG UltraWide 34"', 599.00, '34-inch curved ultrawide monitor with QHD resolution')
ON CONFLICT DO NOTHING;
