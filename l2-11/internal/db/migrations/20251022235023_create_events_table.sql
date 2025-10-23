-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    date DATE NOT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- +goose StatementEnd
