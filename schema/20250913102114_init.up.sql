CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY,
    username      VARCHAR(25) UNIQUE NOT NULL,
    password_hash TEXT               NOT NULL,
    created_at    TIMESTAMP          NOT NULL DEFAULT now(),
    role          SMALLINT           NOT NULL
);

CREATE TABLE accounts
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT    NOT NULL,
    name        TEXT      NOT NULL,
    type        TEXT      NOT NULL,
    is_archived BOOLEAN   NOT NULL,
    order_num   INTEGER   NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE categories
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT    NOT NULL,
    name       TEXT      NOT NULL,
    group_name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE transactions
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT         NOT NULL,
    account_id  BIGINT         NOT NULL,
    date        DATE           NOT NULL,
    amount      NUMERIC(14, 2) NOT NULL,
    category_id BIGINT,
    note        TEXT                    DEFAULT '',
    cleared     BOOLEAN        NOT NULL,
    approved    BOOLEAN        NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP      NOT NULL DEFAULT now()
);

CREATE TABLE budget_allocations
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT         NOT NULL,
    month       INT            NOT NULL,
    year        INT            NOT NULL,
    category_id BIGINT         NOT NULL,
    assigned    NUMERIC(14, 2) NOT NULL,
    created_at  TIMESTAMP      NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP      NOT NULL DEFAULT now()
);

-- Только самые критичные индексы для производительности
CREATE INDEX idx_transactions_user_date ON transactions (user_id, date DESC);
CREATE INDEX idx_transactions_account ON transactions (account_id);
CREATE INDEX idx_budget_user_month ON budget_allocations (user_id, year, month);
