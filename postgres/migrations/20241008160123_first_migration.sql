-- +goose Up
-- +goose StatementBegin
CREATE table auth (
    id serial primary key,
    name text not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table auth;
-- +goose StatementEnd
