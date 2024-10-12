package postgresstorage

import "time"

// StorageInterface интерфейс для работы с бд
type StorageInterface interface {
	CreateChat(users Users) (*ChatID, error)
	DeleteChat(id ChatID) error
	SendMessageChat(message Message) error
}

type (
	// Users пользователи для создания чата
	Users []string
	// ChatID id чата
	ChatID struct {
		ID int64
	}
	// Message сообщение на вставку в чат
	Message struct {
		ChatID ChatID
		From   string
		Body   string
		Time   time.Time
	}
)
