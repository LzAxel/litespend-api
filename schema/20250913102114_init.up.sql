CREATE TABLE users
(
    id            BIGSERIAL PRIMARY KEY UNIQUE,
    username      TEXT UNIQUE,
    password_hash TEXT      NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    role          SMALLINT  NOT NULL
);

CREATE TABLE transaction_categories
(
    id         BIGSERIAL PRIMARY KEY UNIQUE,
    user_id    BIGINT REFERENCES users (id) NOT NULL,
    name       VARCHAR(255) NOT NULL,
    type       SMALLINT     NOT NULL,
    created_at TIMESTAMP    NOT NULL
);

CREATE TABLE goals
(
    id            BIGSERIAL PRIMARY KEY UNIQUE,
    name          VARCHAR(255)     NOT NULL,
    target_amount DOUBLE PRECISION NOT NULL,
    start_amount  DOUBLE PRECISION NOT NULL,
    frequency     SMALLINT         NOT NULL,
    deadline_date TIMESTAMP        NOT NULL,
    created_at    TIMESTAMP        NOT NULL
);

CREATE TABLE transactions
(
    id          BIGSERIAL PRIMARY KEY UNIQUE,
    user_id     BIGINT REFERENCES users (id)                  NOT NULL,
    category_id BIGINT REFERENCES transaction_categories (id) NOT NULL,
    goal_id     BIGINT REFERENCES goals (id),
    description VARCHAR(255)                                  NOT NULL,
    amount      DOUBLE PRECISION                              NOT NULL,
    type        SMALLINT                                      NOT NULL,
    date_time   TIMESTAMP                                     NOT NULL,
    created_at  TIMESTAMP                                     NOT NULL
);

CREATE TABLE prescribed_expanses
(
    id          BIGSERIAL PRIMARY KEY UNIQUE,
    user_id     BIGINT REFERENCES users (id)                  NOT NULL,
    category_id BIGINT REFERENCES transaction_categories (id) NOT NULL,
    description VARCHAR(255)                                  NOT NULL,
    frequency   SMALLINT                                      NOT NULL,
    amount      DOUBLE PRECISION                              NOT NULL,
    date_time   TIMESTAMP                                     NOT NULL,
    created_at  TIMESTAMP                                     NOT NULL
);
