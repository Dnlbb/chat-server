package chatserv

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
)

// SendMessage сервисный хэндлер, вызываем хэндлер базы SendMessageChat.
func (s service) SendMessage(ctx context.Context, message models.Message) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		if errTx = s.storage.SendMessageChat(ctx, message); errTx != nil {
			return fmt.Errorf("error while send message chat: %w", errTx)
		}

		if errTx = s.storage.Log(ctx, models.DELETE); errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("sendMessage error: %w", err)
	}

	return nil
}