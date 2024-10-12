package postgresstorage

import "time"

type storage interface {
	CreateChat(users Users) (*ChatID, error)
	DeleteChat(id ChatID) error
	SendMessageChat(message Message) error
}

type (
	Users  []string
	ChatID struct {
		chatID int64
	}
	Message struct {
		ChatID ChatID
		From   string
		Body   string
		Time   time.Time
	}
)
