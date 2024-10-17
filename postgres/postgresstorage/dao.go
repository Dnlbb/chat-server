package postgresstorage

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"

	sq "github.com/Masterminds/squirrel"
)

// PostgresStorage структура с базой
type PostgresStorage struct {
	con pgx.Conn
}

// InitPostgresStorage инициализация подключения к бд.
func InitPostgresStorage() *PostgresStorage {
	con, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Ошибка при создании подключения к базе данных", err)
	}
	return &PostgresStorage{con: *con}
}

// CloseCon закрываем подключение к бд.
func (s *PostgresStorage) CloseCon() {
	err := s.con.Close(context.Background())
	if err != nil {
		log.Fatal("Ошибка при закрытии соединения с базой данных Postgres", err)
	}
}

// CreateChat создание чата. Для начала идем в таблицу Chats за новым уникальным идентификатором чата.
// Затем добавляем всех переданных пользователей в чат, в случае если мы попытаемся добавить пользователя, которого не
// существует, данный пользователь не будет добавлен и мы перейдем к следующему. Связь пользователей и чатов осуществляется через таблицу
// ChatUsers связью многие ко многим.
func (s *PostgresStorage) CreateChat(users IDs) (*int64, error) {
	var chatID int64

	err := s.con.QueryRow(context.Background(), "INSERT INTO chat_service.chats DEFAULT VALUES RETURNING id").Scan(&chatID)
	if err != nil {
		return nil, fmt.Errorf("error when creating a chat record in the Chats table: %w", err)
	}
	queryBuilder := sq.Insert("chat2users").Columns("chat_id", "user_id")

	for _, userID := range users {
		queryBuilder = queryBuilder.Values(userID, chatID)
	}

	queryBuilder = queryBuilder.Suffix("ON CONFLICT (user_id) DO NOTHING")

	sqlStr, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error when forming the request: %w", err)
	}
	_, err = s.con.Exec(context.Background(), sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("error when executing the request: %w", err)
	}

	log.Printf("Создан чат с id: %d . В него добавлены пользователи: %+v", chatID, users)
	return &chatID, err
}

// DeleteChat удаление чата по id из таблицы Chats
func (s *PostgresStorage) DeleteChat(chatID int64) error {
	_, err := s.con.Exec(context.Background(), "UPDATE chat_service.chats SET is_deleted = TRUE WHERE id = $1", chatID)
	if err != nil {
		return fmt.Errorf("error when deleting a chat: %w", err)
	}

	log.Printf("Чат с id: %d успешно удален", chatID)
	return nil
}

// SendMessageChat отправляем сообщение в определенный чат.
func (s *PostgresStorage) SendMessageChat(message Message) error {
	// Для начала проверяем статус чата (удален или нет), потом уже отправляем в него сообщение.
	var chatIsDeleted bool

	err := s.con.QueryRow(context.Background(), "SELECT is_deleted FROM chat_service.chats WHERE id = $1", message.ChatID).Scan(&chatIsDeleted)
	if err != nil {
		return fmt.Errorf("error checking the chat status: %w", err)
	}

	if chatIsDeleted {
		return fmt.Errorf("error, it is not possible to send a message to a remote chat with an id: %d", message.ChatID)
	}
	var userID int64

	_, err = s.con.Exec(context.Background(), "INSERT INTO chat_service.messages (chat_id, from_user_id, text, timestamp ) VALUES ($1, $2, $3, $4)", message.ChatID, userID, message.Body, message.Time)
	if err != nil {
		log.Println("Ошибка при отправке запроса на добавление сообщения в таблицу Messages")
	}

	log.Printf("Пользователь %s отправил сообщение %s в чат %d", message.From, message.Body, message.ChatID)
	return nil
}
