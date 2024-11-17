package config

import (
	"fmt"
	"net"
	"os"
)

const (
	swaggerHost = "SWAGGER_HOST"
	swaggerPort = "SWAGGER_PORT"
)

// SwaggerConf конфиг для получения адреса для свагера.
type SwaggerConf interface {
	Address() string
}

// SwaggerServerConf структура для реализации интерфейса конфига.т
type SwaggerServerConf struct {
	host string
	port string
}

// NewSwaggerServerConf конструктор для конфига.
func NewSwaggerServerConf() (SwaggerConf, error) {
	host := os.Getenv(swaggerHost)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", swaggerHost)
	}

	port := os.Getenv(swaggerPort)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", swaggerPort)
	}

	return &SwaggerServerConf{
		host: host,
		port: port,
	}, nil
}

// Address получаем адрес свагера.
func (s SwaggerServerConf) Address() string {
	return net.JoinHostPort(s.host, s.port)
}
