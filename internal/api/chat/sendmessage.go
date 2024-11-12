package chat

import (
	"context"
	"fmt"

	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SendMessage конвертация в сервисную модель, затем передача управления в сервисный хэндлер.
func (c *Controller) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*emptypb.Empty, error) {
	message := toMessage(req)
	if err := c.chatService.SendMessage(ctx, *message); err != nil {
		return &emptypb.Empty{}, fmt.Errorf("failed to send message: %w", err)
	}

	return &emptypb.Empty{}, nil
}
