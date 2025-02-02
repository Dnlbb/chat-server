package authrepo

// GetIDs получение id пользователя из Auth по именам.
import (
	"context"
	"fmt"

	authv1 "github.com/Dnlbb/auth/pkg/user_v1"
	"github.com/Dnlbb/chat-server/internal/models"
)

// GetIDs получаем id пользователей из сервиса auth.
func (c AuthRepo) GetIDs(ctx context.Context, usernames models.Usernames) ([]models.ID, error) {
	IDs := make([]models.ID, 0, len(usernames))

	for _, username := range usernames {
		UserProfile, err := c.authClient.Get(ctx, &authv1.GetRequest{NameOrId: &authv1.GetRequest_Username{
			Username: username,
		}})
		if err != nil {
			return nil, fmt.Errorf("error when trying to get a user profile from the authorization service: %w", err)
		}

		IDs = append(IDs, models.ID(UserProfile.Id))
	}

	return IDs, nil
}
