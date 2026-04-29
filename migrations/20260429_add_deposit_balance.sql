ALTER TABLE orders ADD COLUMN deposit decimal(10,2) DEFAULT 0.00;
ALTER TABLE orders ADD COLUMN balance decimal(10,2) DEFAULT 0.00;
ALTER TABLE orders RENAME COLUMN fee TO income;
