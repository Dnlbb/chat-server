package chatserv

import (
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	"github.com/Dnlbb/platform_common/pkg/db"
)

// service сервис.
type service struct {
	storage     repointerface.StorageInterface
	storageAuth repointerface.AuthInterface
	txManager   db.TxManager
}

// NewService инициализация сервиса.
func NewService(storage repointerface.StorageInterface, txManager db.TxManager, storageAuth repointerface.AuthInterface) servinterfaces.ChatService {
	return &service{storage: storage,
		txManager:   txManager,
		storageAuth: storageAuth}
}
