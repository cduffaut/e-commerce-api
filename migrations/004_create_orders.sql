CREATE TABLE IF NOT EXISTS orders (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id),
	status VARCHAR(50) NOT NULL DEFAULT 'pending',
	total_amount NUMERIC(10,2) NOT NULL,
	stripe_payment_id VARCHAR(255),
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

	CONSTRAINT valid_status CHECK (
		status IN ('pending', 'paid', 'cancelled', 'refunded')
	)
);

CREATE TABLE IF NOT EXISTS order_items (
	id BIGSERIAL PRIMARY KEY,
	order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
	product_id BIGINT NOT NULL REFERENCES products(id),
	quantity INTEGER NOT NULL,
	unit_price NUMERIC(10, 2) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status_id ON orders(status);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);