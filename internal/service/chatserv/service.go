package chatserv

import (
	"github.com/Dnlbb/auth/pkg/auth_v1"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	"github.com/Dnlbb/platform_common/pkg/db"
)

// service сервис.
type service struct {
	storage    repointerface.StorageInterface
	txManager  db.TxManager
	authClient auth_v1.AuthClient
}

// NewService инициализация сервиса.
func NewService(storage repointerface.StorageInterface, txManager db.TxManager, authClient auth_v1.AuthClient) servinterfaces.ChatService {
	return &service{storage: storage,
		txManager:  txManager,
		authClient: authClient}
}
