-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
    id  UUID DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL,
    token TEXT NOT NULL,
    updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT user_pk PRIMARY KEY (id),
    CONSTRAINT user_uk_token UNIQUE (token)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "users";
-- +goose StatementEnd
