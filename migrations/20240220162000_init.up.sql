CREATE TABLE IF NOT EXISTS apps
(
    id     SERIAL PRIMARY KEY,
    name   TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    app_id        INT     NOT NULL,
    email         TEXT    NOT NULL UNIQUE,
    password_hash BYTEA   NOT NULL,
    is_admin      BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT fk_app FOREIGN KEY (app_id) REFERENCES apps (id) ON UPDATE CASCADE ON DELETE RESTRICT
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_email ON users (email);
