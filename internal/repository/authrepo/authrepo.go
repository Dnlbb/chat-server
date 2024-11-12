package authrepo

import (
	"github.com/Dnlbb/auth/pkg/auth_v1"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
)

// AuthRepo структура реализующая методы для похода в сервис auth за айдишниками пользователей.
type AuthRepo struct {
	authClient auth_v1.AuthClient
}

// NewAuthRepo конструктор для AuthRepo.
func NewAuthRepo(authClient auth_v1.AuthClient) repointerface.AuthInterface {
	return AuthRepo{authClient: authClient}
}
