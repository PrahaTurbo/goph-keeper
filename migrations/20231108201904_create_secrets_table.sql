-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS secrets (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    type secret_type NOT NULL,
    content BYTEA NOT NULL,
    meta_data BYTEA,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_secrets_user_id ON secrets (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE secrets;
-- +goose StatementEnd
