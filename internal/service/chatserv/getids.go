package chatserv

import (
	"context"
	"fmt"

	authv1 "github.com/Dnlbb/auth/pkg/auth_v1"
	"github.com/Dnlbb/chat-server/internal/models"
)

// GetIDs получение id пользователя из Auth по именам.
func (s service) GetIDs(ctx context.Context, usernames models.Usernames) (models.IDs, error) {
	IDs := models.IDs{}

	for _, username := range usernames {
		UserProfile, err := s.authClient.Get(ctx, &authv1.GetRequest{NameOrId: &authv1.GetRequest_Username{
			Username: username,
		}})
		if err != nil {
			return nil, fmt.Errorf("error when trying to get a user profile from the authorization service: %w", err)
		}

		IDs = append(IDs, UserProfile.Id)
	}

	return IDs, nil
}
