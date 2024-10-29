package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// NewGrpcConfig настраиваем конфиг для запуска сервера, получаем переменные окружения.
func NewGrpcConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if host == "" {
		return nil, fmt.Errorf("error: %w environment variable %s is not set", errors.New("grpc host not found"), grpcHostEnvName)
	}

	port := os.Getenv(grpcPortEnvName)
	if port == "" {
		return nil, fmt.Errorf("error: %w environment variable %s is not set", errors.New("grpc port not found"), grpcPortEnvName)
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

// Address соединяем хост и порт.
func (c *grpcConfig) Address() string {
	return c.host + ":" + c.port
}
