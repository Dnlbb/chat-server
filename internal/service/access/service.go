package access

import (
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
)

// ServiceAccess структура реализующая сервис интерцептора.
type ServiceAccess struct {
	accessRepository repointerface.Access
}

// NewAccessService конструктор для сервиса интерцептора.
func NewAccessService(accessRepository repointerface.Access) servinterfaces.Access {
	return &ServiceAccess{
		accessRepository: accessRepository,
	}
}
