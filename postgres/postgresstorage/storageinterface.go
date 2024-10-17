package postgresstorage

import "time"

// StorageInterface интерфейс для работы с бд
type StorageInterface interface {
	CreateChat(users IDs) (*int64, error)
	DeleteChat(id int64) error
	SendMessageChat(message Message) error
}

type (
	// IDs id пользователей для создания чата
	IDs []int64

	// Message сообщение на вставку в чат
	Message struct {
		ChatID  int64
		From    string
		FromUID int64
		Body    string
		Time    time.Time
	}
)
