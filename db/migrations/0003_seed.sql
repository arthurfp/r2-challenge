-- Seed: Users
INSERT INTO users (id, email, password_hash, name, role)
VALUES
    (gen_random_uuid(), 'admin@example.com', '$2y$10$Y7qkZ/JdFz2H1zYz6bH0hOS8J0bVx1v5o8h3y0yHhJmIu3KqA7nqS', 'Admin', 'admin'),
    (gen_random_uuid(), 'alice@example.com', '$2y$10$Y7qkZ/JdFz2H1zYz6bH0hOS8J0bVx1v5o8h3y0yHhJmIu3KqA7nqS', 'Alice', 'user'),
    (gen_random_uuid(), 'bob@example.com',   '$2y$10$Y7qkZ/JdFz2H1zYz6bH0hOS8J0bVx1v5o8h3y0yHhJmIu3KqA7nqS', 'Bob',   'user')
ON CONFLICT DO NOTHING;

-- Seed: Products
INSERT INTO products (id, name, description, category, price_cents, inventory)
VALUES
    (gen_random_uuid(), 'Laptop 13"', 'Lightweight laptop', 'electronics', 399990, 25),
    (gen_random_uuid(), 'Mechanical Keyboard', 'RGB, blue switches', 'electronics', 12990, 100),
    (gen_random_uuid(), 'Noise-canceling Headphones', 'Over-ear', 'electronics', 25990, 60),
    (gen_random_uuid(), 'Smartphone', '6.1-inch, 128GB', 'electronics', 299990, 40),
    (gen_random_uuid(), 'Coffee Mug', 'Ceramic, 350ml', 'home', 1990, 300),
    (gen_random_uuid(), 'Notebook A5', 'Lined pages', 'stationery', 990, 500),
    (gen_random_uuid(), 'Backpack 20L', 'Water-resistant', 'accessories', 15990, 80),
    (gen_random_uuid(), 'Running Shoes', 'Breathable mesh', 'sports', 22990, 70),
    (gen_random_uuid(), 'Bluetooth Speaker', 'Portable', 'electronics', 8990, 120),
    (gen_random_uuid(), 'LED Desk Lamp', 'Adjustable brightness', 'home', 5990, 150)
ON CONFLICT DO NOTHING;

-- Optional: minimal orders for demo
DO $$
DECLARE u1 UUID; p1 UUID; o1 UUID;
BEGIN
    SELECT id INTO u1 FROM users WHERE email='alice@example.com' LIMIT 1;
    SELECT id INTO p1 FROM products WHERE name='Coffee Mug' LIMIT 1;
    IF u1 IS NOT NULL AND p1 IS NOT NULL THEN
        INSERT INTO orders (id, user_id, status, total_cents)
        VALUES (gen_random_uuid(), u1, 'created', 1990) RETURNING id INTO o1;
        INSERT INTO order_items (order_id, product_id, quantity, price_cents)
        VALUES (o1, p1, 1, 1990);
    END IF;
END $$;


