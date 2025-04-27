CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  customer_id VARCHAR(255) NOT NULL,
  status INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_orders_user_id ON orders(customer_id);
