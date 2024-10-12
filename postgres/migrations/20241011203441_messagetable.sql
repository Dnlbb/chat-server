-- +goose Up
-- +goose StatementBegin
CREATE TABLE Users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE Chats (
                       id SERIAL PRIMARY KEY,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       is_deleted BOOLEAN DEFAULT FALSE
);

CREATE TABLE ChatUsers (
                           id SERIAL PRIMARY KEY,
                           chat_id INTEGER REFERENCES Chats(id) ON DELETE CASCADE,
                           user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE,
                           joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           UNIQUE(chat_id, user_id)
);

CREATE TABLE Messages (
                          id SERIAL PRIMARY KEY,
                          chat_id INTEGER REFERENCES Chats(id) ON DELETE CASCADE,
                          from_user_id INTEGER REFERENCES Users(id) ON DELETE SET NULL,
                          text TEXT NOT NULL,
                          timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON Users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_chats_updated_at
    BEFORE UPDATE ON Chats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Messages;
DROP TABLE IF EXISTS ChatUsers;
DROP TABLE IF EXISTS Chats;
DROP TABLE IF EXISTS Users;

DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
-- +goose StatementEnd
