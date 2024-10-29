package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

type pgConfig struct {
	dsn string
}

// NewPgConfig настраиваем конфиг для запуска базы данных и получаем значения переменных окружения.
func NewPgConfig() (PGConfig, error) {
	DSN := os.Getenv(dsnEnvName)
	if DSN == "" {
		return nil, fmt.Errorf("error: %w environment variable %s not set", errors.New("dsn env not found"), dsnEnvName)
	}

	return &pgConfig{
		dsn: DSN,
	}, nil
}

// DSN реализуем интерфейс и возвращаем dsn.
func (c *pgConfig) DSN() string {
	return c.dsn
}
