-- +goose Up
-- +goose StatementBegin
CREATE TYPE secret_type AS ENUM (
    'CREDENTIALS',
    'TEXT',
    'BINARY',
    'CARD'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE secret_type;
-- +goose StatementEnd