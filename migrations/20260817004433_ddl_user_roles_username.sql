-- +goose Up
CREATE TYPE lang AS ENUM ('eng', 'fa');

CREATE TYPE role_name AS ENUM ('admin', 'owner', 'editor', 'viewer');

CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name role_name NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE user_roles (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    role_id INTEGER REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE (user_id, role_id)
);

CREATE TABLE usernames (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE ON UPDATE CASCADE,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255),
    is_verified_email BOOLEAN DEFAULT FALSE,
    day_duration INTEGER NULL,
    period INTEGER NULL,
    user_lang lang NOT NULL,
    description TEXT,
    UNIQUE(id, username),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_usernames_user_id ON usernames(user_id);

CREATE INDEX idx_usernames_role_id ON usernames(role_id);

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);

CREATE INDEX idx_roles_deleted_at ON roles(deleted_at);

CREATE INDEX idx_usernames_deleted_at ON usernames(deleted_at);

-- +goose Down
DROP INDEX IF EXISTS idx_usernames_deleted_at;

DROP INDEX IF EXISTS idx_roles_deleted_at;

DROP INDEX IF EXISTS idx_users_deleted_at;

DROP INDEX IF EXISTS idx_usernames_user_id;

DROP INDEX IF EXISTS idx_usernames_role_id;

DROP INDEX IF EXISTS idx_user_roles_user_id;

DROP TABLE IF EXISTS usernames;

DROP TABLE IF EXISTS user_roles;

DROP TABLE IF EXISTS roles;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS role_name;

DROP TYPE IF EXISTS lang;