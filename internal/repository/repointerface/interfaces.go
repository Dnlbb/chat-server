package repointerface

import (
	"context"

	"github.com/Dnlbb/chat-server/internal/models"
)

// StorageInterface интерфейс для работы с бд
type StorageInterface interface {
	CreateChat(ctx context.Context, id models.IDs) (*int64, error)
	DeleteChat(ctx context.Context, chatID models.Chat) error
	SendMessageChat(ctx context.Context, message models.Message) error
	Log(ctx context.Context, key models.LogKey) error
}

// AuthInterface интерфейс для получения ID пользователей из сервиса auth.
type AuthInterface interface {
	GetIDs(ctx context.Context, usernames models.Usernames) ([]models.ID, error)
}
