package storage

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
	"github.com/Dnlbb/platform_common/pkg/db"
)

// DeleteChat удаление чата по id из таблицы Chats
func (s *storage) DeleteChat(ctx context.Context, chat models.Chat) error {
	q := db.Query{
		Name:     "Delete chat",
		QueryRow: "UPDATE chat_service.chats SET is_deleted = TRUE WHERE id = $1",
	}

	_, err := s.db.DB().ExecContext(ctx, q, chat.ID)
	if err != nil {
		return fmt.Errorf("error when deleting a chat: %w", err)
	}

	return nil
}
