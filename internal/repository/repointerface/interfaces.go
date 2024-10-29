package repointerface

import (
	"context"

	"github.com/Dnlbb/chat-server/internal/models"
)

// StorageInterface интерфейс для работы с бд
type StorageInterface interface {
	CreateChat(ctx context.Context, id models.IDs) (*int64, error)
	DeleteChat(ctx context.Context, chatID models.ChatID) error
	SendMessageChat(ctx context.Context, message models.Message) error
	Log(ctx context.Context, key models.LogKey) error
}
