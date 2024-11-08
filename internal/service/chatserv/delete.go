package chatserv

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
)

// Delete сервисный хэндлер, вызываем хэндлер базы DeleteChat.
func (s service) Delete(ctx context.Context, chatID models.ChatID) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		if errTx = s.storage.DeleteChat(ctx, chatID); errTx != nil {
			return fmt.Errorf("error delete chat: %w", errTx)
		}

		if errTx = s.storage.Log(ctx, models.DELETE); errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("chat delete error: %w", err)
	}

	return nil
}
