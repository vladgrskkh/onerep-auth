-- +goose Up
CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    display_name TEXT NOT NULL,
    avatar_url TEXT,
    gender TEXT NOT NULL DEFAULT 'other',
    birth_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_users_email ON auth.users (email);

-- +goose Down
DROP TABLE IF EXISTS auth.users;
DROP SCHEMA IF EXISTS auth;
