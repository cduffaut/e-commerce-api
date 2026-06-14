CREATE TABLE IF NOT EXISTS products (
	id BIGSERIAL PRIMARY KEY,
	name	VARCHAR(255) NOT NULL,
	description TEXT,
	price NUMERIC(10, 2) NOT NULL,
	stock INT NOT NULL DEFAULT 0,
	category VARCHAR(100),
	is_active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
