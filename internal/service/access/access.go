package access

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	authPrefix = "Bearer "
)

// Access метод дял проверки доступа пользователя.
func (as ServiceAccess) Access(ctx context.Context, path string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("metadata error")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return fmt.Errorf("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return fmt.Errorf("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)
	mod := metadata.New(map[string]string{"Authorization": "Bearer " + accessToken})
	clientCtx := metadata.NewOutgoingContext(ctx, mod)

	return as.accessRepository.Access(clientCtx, path)
}
