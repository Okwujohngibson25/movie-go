-- +goose Up
-- +goose StatementBegin

--create custom type
CREATE TYPE account_role AS ENUM ('user', 'admin');

--create user table
CREATE TABLE users(
    id UUID PRIMARY KEY,
    fullname TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    role account_role NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS account_role;


-- +goose StatementEnd
