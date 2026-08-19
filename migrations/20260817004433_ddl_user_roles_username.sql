-- +goose Up
CREATE TYPE lang AS ENUM ('eng', 'fa');

CREATE TYPE role_name AS ENUM ('admin', 'owner', 'editor', 'viewer');

CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    user_lang lang NOT NULL DEFAULT 'fa',
    -- Telegram User ID 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE target_accounts (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    day_duration INTEGER,
    period INTEGER,
    description TEXT,
    enable BOOLEAN DEFAULT false,
    UNIQUE(owner_user_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- HACK : 
-- CREATE TABLE casbin_rule (
--     id BIGSERIAL PRIMARY KEY,
--     ptype VARCHAR(100) NOT NULL DEFAULT 'g',
--     -- Subject : matin
--     v0 VARCHAR(100),
--     -- Role: owner
--     v1 VARCHAR(100),
--     -- Resource: karim-account
--     v2 VARCHAR(100),
--     v3 VARCHAR(100),
--     v4 VARCHAR(100),
--     v5 VARCHAR(100),
--     CONSTRAINT unique_casbin_rule UNIQUE (ptype, v0, v1, v2, v3, v4, v5)
-- );
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_targets_owner ON target_accounts(owner_user_id);

CREATE INDEX idx_targets_deleted_at ON target_accounts(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS target_accounts;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS role_name;

DROP TYPE IF EXISTS lang;