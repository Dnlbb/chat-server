package storage

import (
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/platform_common/pkg/db"
)

type storage struct {
	db db.Client
}

// NewPostgresRepo инициализируем хранилище postgresql и приводим его к типу интерфейса StorageInterface.
func NewPostgresRepo(db db.Client) repointerface.StorageInterface {
	return &storage{db: db}
}
