package config

import (
	"net"
	"os"

	"github.com/pkg/errors"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

// HTTPConfig конфиг для получения адреса HTTP сервера.
type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig конструктор дял конфига.
func NewHTTPConfig() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.Errorf("%s must be set", httpHostEnvName)
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.Errorf("%s must be set", httpPortEnvName)
	}

	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
