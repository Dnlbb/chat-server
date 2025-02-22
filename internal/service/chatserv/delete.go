package chatserv

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
)

// Delete сервисный хэндлер, вызываем хэндлер базы DeleteChat.
func (s service) Delete(ctx context.Context, chat models.Chat) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		if errTx = s.storage.DeleteChat(ctx, chat); errTx != nil {
			return fmt.Errorf("error delete chat: %w", errTx)
		}

		if errTx = s.storage.Log(ctx, models.DELETE); errTx != nil {
			return fmt.Errorf("error logging delete chat: %w", errTx)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
