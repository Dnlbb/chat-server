package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/Dnlbb/chat-server/internal/models"
	"github.com/Dnlbb/platform_common/pkg/db"
)

// SendMessageChat отправляем сообщение в определенный чат.
func (s *storage) SendMessageChat(ctx context.Context, message models.Message) error {
	// Для начала проверяем статус чата (удален или нет), потом уже отправляем в него сообщение.
	var chatIsDeleted models.IsDeleted

	q := db.Query{
		Name:     "Check status chat",
		QueryRow: "SELECT is_deleted FROM chat_service.chats WHERE id = $1",
	}

	if err := s.db.DB().ScanOneContext(ctx, &chatIsDeleted, q); err != nil {
		return fmt.Errorf("error checking the chat status: %w", err)
	}

	if chatIsDeleted.Flag {
		return fmt.Errorf("error, it is not possible to send a message to a remote chat with an id: %d", message.ChatID)
	}

	q = db.Query{
		Name:     "Add message to chat",
		QueryRow: "INSERT INTO chat_service.messages (chat_id, from_user_id, text, timestamp ) VALUES ($1, $2, $3, $4)",
	}

	_, err := s.db.DB().ExecContext(ctx, q, message.ChatID, message.FromUID, message.Body, message.Time)
	if err != nil {
		log.Println("Ошибка при отправке запроса на добавление сообщения в таблицу Messages")
	}

	return nil
}
