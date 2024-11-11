package chatserv

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
)

// Create сервисный хэндлер, вызываем хэндлер базы CreateChat.
func (s service) Create(ctx context.Context, usernames models.Usernames) (*int64, error) {
	var ChatID *int64

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		IDs, errTx := s.storageAuth.GetIDs(ctx, usernames)
		if errTx != nil {
			return fmt.Errorf("error with get IDs: %w", errTx)
		}

		ChatID, errTx = s.storage.CreateChat(ctx, IDs)
		if errTx != nil {
			return fmt.Errorf("error creating chat: %w", errTx)
		}

		errTx = s.storage.Log(ctx, models.CREATE)
		if errTx != nil {
			return fmt.Errorf("error logging create chat: %w", errTx)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return ChatID, nil
}
