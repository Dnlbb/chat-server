package chat

import (
	"github.com/Dnlbb/auth/pkg/auth_v1"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
)

// Controller структура реализующая сгенерированный grpc сервер chat.
type Controller struct {
	chatv1.UnimplementedChatServer
	chatService servinterfaces.ChatService
	authClient  auth_v1.AuthClient
}

// NewController конструктор для реализации grpc сервера chat.
func NewController(chatService servinterfaces.ChatService, authClient auth_v1.AuthClient) *Controller {
	return &Controller{chatService: chatService, authClient: authClient}
}
