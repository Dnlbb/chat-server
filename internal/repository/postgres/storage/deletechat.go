package storage

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/client/db"
	"github.com/Dnlbb/chat-server/internal/models"
)

// DeleteChat удаление чата по id из таблицы Chats
func (s *storage) DeleteChat(ctx context.Context, chatID models.ChatID) error {
	q := db.Query{
		Name:     "Delete chat",
		QueryRow: "UPDATE chat_service.chats SET is_deleted = TRUE WHERE id = $1",
	}

	_, err := s.db.DB().ExecContext(ctx, q, chatID.ID)
	if err != nil {
		return fmt.Errorf("error when deleting a chat: %w", err)
	}

	return nil
}
