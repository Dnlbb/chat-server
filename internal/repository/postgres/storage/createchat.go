package storage

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
	"github.com/Dnlbb/platform_common/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

// CreateChat создание чата. Для начала идем в таблицу Chats за новым уникальным идентификатором чата.
// Затем добавляем всех переданных пользователей в чат, в случае если мы попытаемся добавить пользователя, которого не
// существует, данный пользователь не будет добавлен и мы перейдем к следующему. Связь пользователей и чатов осуществляется через таблицу
// ChatUsers связью многие ко многим.
func (s *storage) CreateChat(ctx context.Context, IDs models.IDs) (*int64, error) {
	var chatID models.ChatID

	q := db.Query{
		Name:     "Get chatID",
		QueryRow: "INSERT INTO chat_service.chats DEFAULT VALUES RETURNING id",
	}

	if err := s.db.DB().ScanOneContext(ctx, &chatID.ID, q); err != nil {
		return nil, fmt.Errorf("error when creating a chat record in the Chats table: %w", err)
	}

	queryBuilder := sq.Insert("chat_service.chat2users").Columns("chat_id", "user_id")

	for _, userID := range IDs {
		queryBuilder = queryBuilder.Values(chatID.ID, userID)
	}

	queryBuilder = queryBuilder.PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error when forming the request: %w", err)
	}

	q = db.Query{
		Name:     "Mapping users to chats",
		QueryRow: sqlStr,
	}

	_, err = s.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("error when executing the request: %w", err)
	}

	return &chatID.ID, err
}
