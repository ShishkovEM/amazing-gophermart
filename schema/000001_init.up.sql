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

CREATE TABLE IF NOT EXISTS operations (
    id UUID,
    user_id UUID,
    operation_type VARCHAR(20) NOT NULL,
    order_num VARCHAR(255),
    amount FLOAT4,
    status VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    CONSTRAINT operations_pkey PRIMARY KEY (id),
    CONSTRAINT operations_users_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);

CREATE OR REPLACE VIEW balance AS (
  WITH total AS (
    SELECT
      user_id,
      SUM(
        CASE WHEN operation_type = 'order' THEN amount ELSE 0 END
      ) total
    FROM
      operations
    GROUP BY
      user_id
  ),
  withdraw AS (
    SELECT
      user_id,
      SUM(
        CASE WHEN operation_type = 'withdraw' THEN amount ELSE 0 END
      ) withdraw
    FROM
      operations
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

DELETE FROM operations WHERE user_id IN (SELECT user_id FROM users WHERE username LIKE 'test%');

DELETE FROM users WHERE username LIKE 'test%';