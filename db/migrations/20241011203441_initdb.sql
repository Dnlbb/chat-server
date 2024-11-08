-- +goose Up
-- +goose StatementBegin

CREATE SCHEMA IF NOT EXISTS chat_service;

CREATE TABLE chat_service.chats (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE chat_service.chat2users (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT REFERENCES chat_service.chats(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(chat_id, user_id)
);

CREATE TABLE chat_service.messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT REFERENCES chat_service.chats(id) ON DELETE CASCADE,
    from_user_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS chat_service.messages;
DROP TABLE IF EXISTS chat_service.chat2users;
DROP TABLE IF EXISTS chat_service.chats;
DROP SCHEMA IF EXISTS chat_service;

-- +goose StatementEnd

