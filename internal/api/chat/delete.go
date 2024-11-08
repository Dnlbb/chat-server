package chat

import (
	"context"
	"fmt"

	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete конвертация в сервисную модель chatID, а потом передача управления в сервисный хэндлер.
func (c *Controller) Delete(ctx context.Context, req *chatv1.DeleteRequest) (*emptypb.Empty, error) {
	chatID := toModelsChatID(req)

	if err := c.chatService.Delete(ctx, chatID); err != nil {
		return &emptypb.Empty{}, fmt.Errorf("delete chat error: %w", err)
	}

	return &emptypb.Empty{}, nil
}
