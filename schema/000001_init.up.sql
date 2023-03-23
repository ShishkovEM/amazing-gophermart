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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT orders_pkey PRIMARY KEY (order_num),
    CONSTRAINT orders_users_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS withdrawals (
    user_id UUID,
    order_num VARCHAR(255),
    withdraw FLOAT4,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT withdrawals_pkey PRIMARY KEY (order_num),
    CONSTRAINT withdrawals_users_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);

DELETE FROM withdrawals WHERE user_id IN (SELECT user_id FROM users WHERE username LIKE 'test%');

DELETE FROM orders WHERE user_id IN (SELECT user_id FROM users WHERE username LIKE 'test%');

DELETE FROM users WHERE username LIKE 'test%';

CREATE
OR REPLACE VIEW balance AS (
  WITH total AS (
    SELECT
      user_id,
      SUM(accrual) total
    FROM
      orders
    GROUP BY
      user_id
  ),
  withdraw AS (
    SELECT
      user_id,
      SUM(withdraw) withdraw
    FROM
      withdrawals
    GROUP BY
      user_id
  )
  SELECT
    users.id user_id,
    COALESCE(total.total, 0):: NUMERIC(10, 2) total,
    COALESCE(withdraw.withdraw, 0):: NUMERIC(10, 2) withdraw,
    (COALESCE(total.total, 0) - COALESCE(withdraw.withdraw, 0)):: NUMERIC(10, 2) "current"
  FROM
    users
    LEFT JOIN total ON total.user_id = users.id
    LEFT JOIN withdraw ON withdraw.user_id = users.id
);