-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    description TEXT DEFAULT '',
    user_id UUID NOT NULL,
    reminder TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
