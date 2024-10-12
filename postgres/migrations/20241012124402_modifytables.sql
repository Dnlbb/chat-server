-- +goose Up
-- +goose StatementBegin

ALTER TABLE Users
    DROP COLUMN updated_at;


ALTER TABLE Chats
    DROP COLUMN updated_at;


DROP TRIGGER IF EXISTS update_users_updated_at ON Users;
DROP TRIGGER IF EXISTS update_chats_updated_at ON Chats;


DROP FUNCTION IF EXISTS update_updated_at_column;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE Users
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;


ALTER TABLE Chats
    ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;


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
