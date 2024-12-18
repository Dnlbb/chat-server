package models

import (
	"time"
)

// LogKey тип для логирования запросов в базу.
type LogKey string

// Ключи для логирования запросов в базу.
const (
	CREATE      LogKey = "create"
	DELETE      LogKey = "delete"
	SENDMESSAGE LogKey = "sendMessage"
)

type (
	// IDs id пользователей для создания чата.
	IDs []ID

	// ID id пользователя.
	ID int64

	// Chat id чата.
	Chat struct {
		ID int64 `json:"id"`
	}

	// IsDeleted флаг проверки чата.
	IsDeleted struct {
		Flag bool `db:"is_deleted"`
	}

	// Message сообщение на вставку в чат.
	Message struct {
		ChatID    int64
		FromUname string
		FromUID   int64
		Body      string
		Time      time.Time
	}
	// Usernames имена пользователей.
	Usernames []string
)
