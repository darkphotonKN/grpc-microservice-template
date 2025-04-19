CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  price_id VARCHAR(255) NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id)
);

