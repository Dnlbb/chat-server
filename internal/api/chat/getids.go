package chat

import (
	"context"
	"fmt"
	"log"

	authv1 "github.com/Dnlbb/auth/pkg/auth_v1"
	"github.com/Dnlbb/chat-server/internal/models"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
)

// GetIDs получаем id пользователей для создания чата из сервиса авторизации.
func (c *Controller) GetIDs(ctx context.Context, req *chatv1.CreateRequest) (models.IDs, error) {
	IDs := models.IDs{}

	for _, username := range req.Usernames {
		UserProfile, err := c.authClient.Get(ctx, &authv1.GetRequest{NameOrId: &authv1.GetRequest_Username{
			Username: username,
		}})
		if err != nil {
			return nil, fmt.Errorf("error when trying to get a user profile from the authorization service: %w", err)
		}

		IDs = append(IDs, UserProfile.Id)

		log.Printf("Получен профиль пользователя %s: %v", username, UserProfile)
	}

	return IDs, nil
}
