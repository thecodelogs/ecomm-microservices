CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id              UUID         PRIMARY KEY DEFAULT uuidv7(),
    role_id         UUID         NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    email           VARCHAR(255) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    first_name      VARCHAR(100) NOT NULL,
    last_name       VARCHAR(100) NOT NULL,
    phone           VARCHAR(20),
    avatar_url      TEXT,
    email_verified  BOOLEAN      NOT NULL DEFAULT FALSE,
    phone_verified  BOOLEAN      NOT NULL DEFAULT FALSE,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email    ON users (email);
CREATE INDEX idx_users_role_id  ON users (role_id);
CREATE INDEX idx_users_active   ON users (is_active) WHERE is_active = TRUE;
