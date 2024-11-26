package access

import (
	"context"

	"github.com/Dnlbb/auth/pkg/auth_v1"
)

// Access вызываем сервис авторизации для проверки доступа.
func (repo RepoAccess) Access(ctx context.Context, path string) error {
	_, err := repo.client.Check(ctx, &auth_v1.CheckRequest{EndpointAddress: path})
	return err
}
