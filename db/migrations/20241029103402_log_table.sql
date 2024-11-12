-- +goose Up
-- +goose StatementBegin
CREATE TABLE log (
    id BIGSERIAL primary key,
    name TEXT,
    time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE log;
-- +goose StatementEnd
