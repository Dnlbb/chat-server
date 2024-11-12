package storage

import (
	"context"
	"fmt"

	"github.com/Dnlbb/chat-server/internal/models"
	"github.com/Dnlbb/platform_common/pkg/db"
	sq "github.com/Masterminds/squirrel"
)

func (s *storage) Log(ctx context.Context, key models.LogKey) error {
	query := sq.Insert("log").Columns("name")

	query = query.Values(key)

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
