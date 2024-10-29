package storage

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/client/db"
	"github.com/Dnlbb/chat-server/internal/models"
	sq "github.com/Masterminds/squirrel"
)

func (s *storage) Log(ctx context.Context, key models.LogKey) error {
	query := sq.Insert("log").Columns("name")

	switch key {
	case models.CREATE:
		query = query.Values(models.CREATE)
	case models.DELETE:
		query = query.Values(models.DELETE)
	case models.SENDMESSAGE:
		query = query.Values(models.SENDMESSAGE)
	}

	query = query.PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("error with make sql query: %w", err)
	}

	q := db.Query{
		Name:     "Log",
		QueryRow: sql,
	}

	_, err = s.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error with insert to logging table: %w", err)
	}

	return nil
}
