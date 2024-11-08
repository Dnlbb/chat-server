package servinterfaces

import (
	"context"

	"github.com/Dnlbb/chat-server/internal/models"
)

// ChatService интерфейс чат сервиса.
type ChatService interface {
	Create(ctx context.Context, IDs models.IDs) (*int64, error)
	Delete(ctx context.Context, chatID models.ChatID) error
	SendMessage(ctx context.Context, message models.Message) error
}
