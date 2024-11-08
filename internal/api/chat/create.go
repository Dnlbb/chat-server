package chat

import (
	"context"
	"fmt"

	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
)

// Create получение списка id пользователей для добавления в чат и сервиса авторизации
// и дальнейшая их передача в сервисный слой.
func (c *Controller) Create(ctx context.Context, req *chatv1.CreateRequest) (*chatv1.CreateResponse, error) {
	IDs, err := c.chatService.GetIDs(ctx, toUsernamesModels(req))
	if err != nil {
		return nil, fmt.Errorf("error with get IDs: %w", err)
	}

	chatID, err := c.chatService.Create(ctx, IDs)
	if err != nil {
		return nil, fmt.Errorf("error when trying to create a chat: %w", err)
	}

	return &chatv1.CreateResponse{Id: *chatID}, nil
}
