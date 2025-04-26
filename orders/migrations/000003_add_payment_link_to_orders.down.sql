-- Down migration to remove the payment_link column
ALTER TABLE orders DROP COLUMN payment_link;
