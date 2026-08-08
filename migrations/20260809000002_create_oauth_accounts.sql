-- +goose Up
CREATE TABLE auth.oauth_accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_user_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE(provider, provider_user_id)
);

CREATE INDEX idx_oauth_accounts_user_id ON auth.oauth_accounts (user_id);

-- +goose Down
DROP TABLE IF EXISTS auth.oauth_accounts;
