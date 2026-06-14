CREATE TABLE IF NOT EXISTS cart_items (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
	quantity INTEGER NOT NULL DEFAULT 1,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

	CONSTRAINT unique_user_product UNIQUE (user_id, product_id),
	CONSTRAINT quantity_positive CHECK (quantity > 0)
);

CREATE INDEX IF NOT EXISTS idx_cart_items_user_id ON cart_items(user_id);