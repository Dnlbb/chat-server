package chatserv

import (
	"github.com/Dnlbb/chat-server/internal/client/db"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
)

// service сервис.
type service struct {
	storage   repointerface.StorageInterface
	txManager db.TxManager
}

// NewService инициализация сервиса.
func NewService(storage repointerface.StorageInterface, txManager db.TxManager) servinterfaces.ChatService {
	return &service{storage: storage,
		txManager: txManager}
}
