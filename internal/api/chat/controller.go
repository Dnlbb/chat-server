package chat

import (
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
)

// Controller структура реализующая сгенерированный grpc сервер chat.
type Controller struct {
	chatv1.UnimplementedChatServer
	chatService servinterfaces.ChatService
}

// NewController конструктор для реализации grpc сервера chat.
func NewController(chatService servinterfaces.ChatService) *Controller {
	return &Controller{chatService: chatService}
}
