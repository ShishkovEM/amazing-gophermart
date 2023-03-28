CREATE TABLE IF NOT EXISTS users (
    id UUID,
    username VARCHAR(255),
    pass VARCHAR(255),
    token TEXT,
    token_expires TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT users_pkey PRIMARY KEY (id),
    CONSTRAINT username_unique UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS orders (
    user_id UUID,
    order_num VARCHAR(255),
    accrual FLOAT4,
    status VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT orders_pkey PRIMARY KEY (order_num),
    CONSTRAINT orders_users_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);

DELETE FROM orders WHERE user_id IN (SELECT user_id FROM users WHERE username LIKE 'test%');

DELETE FROM users WHERE username LIKE 'test%';