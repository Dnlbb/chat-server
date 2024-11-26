package authrepo

import (
	"github.com/Dnlbb/auth/pkg/user_v1"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
)

// AuthRepo структура реализующая методы для похода в сервис auth за айдишниками пользователей.
type AuthRepo struct {
	authClient user_v1.UserApiClient
}

// NewAuthRepo конструктор для AuthRepo.
func NewAuthRepo(authClient user_v1.UserApiClient) repointerface.AuthInterface {
	return AuthRepo{authClient: authClient}
}
