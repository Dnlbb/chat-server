package servinterfaces

import (
	"context"

	"github.com/Dnlbb/chat-server/internal/models"
)

// ChatService интерфейс чат сервиса.
type ChatService interface {
	Create(ctx context.Context, usernames models.Usernames) (*int64, error)
	Delete(ctx context.Context, chat models.Chat) error
	SendMessage(ctx context.Context, message models.Message) error
}

// Access интерфейс реализующий проверку.
type Access interface {
	Access(ctx context.Context, path string) error
}
