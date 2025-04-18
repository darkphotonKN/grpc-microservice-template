CREATE TABLE order_items (
  id VARCHAR(36) PRIMARY KEY,
  order_id VARCHAR(36) NOT NULL,
  name VARCHAR(255) NOT NULL,
  quantity INT NOT NULL,
  price_id VARCHAR(255) NOT NULL,
  FOREIGN KEY (order_id) REFERENCES orders(id)
);

