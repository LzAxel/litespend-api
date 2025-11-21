CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY UNIQUE,
    username      VARCHAR(25) UNIQUE,
    password_hash TEXT      NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    role          SMALLINT  NOT NULL
);

CREATE TABLE transaction_categories
(
    id         BIGSERIAL PRIMARY KEY UNIQUE,
    user_id    BIGINT REFERENCES users (id) NOT NULL,
    name       VARCHAR(255)                 NOT NULL,
    type       SMALLINT                     NOT NULL CHECK (type IN (0, 1)) NOT NULL,
    created_at TIMESTAMP                    NOT NULL,
    UNIQUE (user_id, name, type)
);

CREATE TABLE recurring_bills
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name        TEXT           NOT NULL,
    category_id BIGINT REFERENCES transaction_categories (id),
    amount      NUMERIC(14, 2) NOT NULL,
    day_due     SMALLINT CHECK (day_due BETWEEN 1 AND 31),
    is_active   BOOLEAN,
    created_at  TIMESTAMP
);

CREATE TABLE bill_instances
(
    id                BIGSERIAL PRIMARY KEY,
    recurring_bill_id BIGINT         NOT NULL REFERENCES recurring_bills (id),
    year              INT            NOT NULL,
    month             INT            NOT NULL CHECK (month BETWEEN 1 AND 12),
    amount_expected   NUMERIC(14, 2) NOT NULL,
    created_at        TIMESTAMP,
    UNIQUE (recurring_bill_id, year, month)
);

CREATE TABLE transactions
(
    id               BIGSERIAL PRIMARY KEY UNIQUE,
    user_id          BIGINT REFERENCES users (id) NOT NULL,
    category_id      BIGINT REFERENCES transaction_categories (id),
    bill_instance_id BIGINT REFERENCES bill_instances (id),
    description      VARCHAR(255)                 NOT NULL,
    amount           NUMERIC(14, 2)               NOT NULL,
    date             TIMESTAMP                    NOT NULL,
    created_at       TIMESTAMP                    NOT NULL
);

CREATE TABLE budgets
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT         NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    category_id BIGINT         NOT NULL REFERENCES transaction_categories (id),
    year        INT            NOT NULL,
    month       INT            NOT NULL CHECK (month BETWEEN 1 AND 12),
    budgeted    NUMERIC(14, 2) NOT NULL,
    created_at  TIMESTAMP,
    UNIQUE (user_id, category_id, year, month)
);