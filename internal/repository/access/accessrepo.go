package access

import (
	"github.com/Dnlbb/auth/pkg/auth_v1"
)

// RepoAccess структура реализующая доступ к клиенту для интерцептора.
type RepoAccess struct {
	client auth_v1.AuthClient
}

// NewAccessRepo конструктор для структуры реализующей доступ к клиенту для интерцептора.
func NewAccessRepo(client auth_v1.AuthClient) *RepoAccess {
	return &RepoAccess{client: client}
}
