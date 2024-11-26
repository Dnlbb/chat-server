package access

import (
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
)

// AuthInterceptor структура реализующая интерцептор для авторизации.
type AuthInterceptor struct {
	AccessService servinterfaces.Access
}

// NewAuthInterceptor конструктор для структуры реализующей интерцептор для авторизации.
func NewAuthInterceptor(accessService servinterfaces.Access) *AuthInterceptor {
	return &AuthInterceptor{
		AccessService: accessService,
	}
}
