package postgresstorage

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

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
func (s *PostgresStorage) CreateChat(users Users) (*ChatID, error) {
	var chatID ChatID

	err := s.con.QueryRow(context.Background(), "INSERT INTO Chats DEFAULT VALUES RETURNING id").Scan(&chatID)
	if err != nil {
		log.Println("Ошибка при создании записи чата в таблице Chats", err)
		return nil, err
	}
	queryBuilder := sq.Insert("ChatUsers").Columns("chat_id", "user_id")

	for _, username := range users {
		var userID int64
		err := s.con.QueryRow(context.Background(), "SELECT id FROM Users WHERE username = $1", username).Scan(&userID)
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("Пользователя которого пытаются добавить в чат не существует", err)
		} else if err != nil {
			log.Println("Ошибка при попытке получить id пользователя из таблицы Users")
			return nil, err
		}
		queryBuilder = queryBuilder.Values(userID, chatID)
	}

	queryBuilder = queryBuilder.Suffix("ON CONFLICT (user_id) DO NOTHING")

	sqlStr, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка при формировании запроса: %w", err)
	}
	_, err = s.con.Exec(context.Background(), sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}

	log.Printf("Создан чат с id: %d . В него добавлены пользователи: %+v", chatID, users)
	return &chatID, err
}

// DeleteChat удаление чата по id из таблицы Chats
func (s *PostgresStorage) DeleteChat(chatID ChatID) error {
	_, err := s.con.Exec(context.Background(), "UPDATE Chats SET is_deleted = TRUE WHERE id = $1", chatID)
	if err != nil {
		log.Println("Ошибка при удалении чата:", err)
		return err
	}

	log.Printf("Чат с id: %d успешно удален", chatID)
	return nil
}

// SendMessageChat отправляем сообщение в определенный чат.
func (s *PostgresStorage) SendMessageChat(message Message) error {
	// Для начала проверяем статус чата (удален или нет), потом уже отправляем в него сообщение.
	var chatIsDeleted bool

	err := s.con.QueryRow(context.Background(), "SELECT is_deleted FROM Chats WHERE id = $1", message.ChatID).Scan(&chatIsDeleted)
	if err != nil {
		log.Println("Ошибка при проверке состояния чата", err)
		return err
	}

	if chatIsDeleted {
		return fmt.Errorf("невозможно отправить сообщение в удалённый чат с id: %d", message.ChatID)
	}
	var userID int64

	err = s.con.QueryRow(context.Background(), "SELECT id FROM Users WHERE username = $1", message.From).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Println("Пользователь, который пытается написать сообщение не существует", err)
	} else if err != nil {
		log.Println("Ошибка при попытке получения id пользователя из бд")
		return err
	}
	_, err = s.con.Exec(context.Background(), "INSERT INTO Messages (chat_id, from_user_id, text, timestamp ) VALUES ($1, $2, $3, $4)", message.ChatID, userID, message.Body, message.Time)
	if err != nil {
		log.Println("Ошибка при отправке запроса на добавление сообщения в таблицу Messages")
	}

	log.Printf("Пользователь %s отправил сообщение %s в чат %d", message.From, message.Body, message.ChatID)
	return nil
}
