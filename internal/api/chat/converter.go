package chat

import (
	"time"

	"github.com/Dnlbb/chat-server/internal/models"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toModelsChatID конвертация в сервисную модель ChatID.
func toModelsChatID(req *chatv1.DeleteRequest) models.ChatID {
	return models.ChatID{
		ID: req.Id,
	}
}

// toModelsMessage конвертация в сервисную модель Message.
func toModelsMessage(req *chatv1.SendMessageRequest) *models.Message {
	return &models.Message{
		ChatID:    req.GetChatID(),
		FromUname: req.GetFromUname(),
		FromUID:   req.GetFromUid(),
		Body:      req.GetBody(),
		Time:      toTimestampTime(req.GetTime()),
	}
}

// ToTimestampProto конвертация proto времени в обычный time.Time.
func toTimestampTime(time *timestamppb.Timestamp) time.Time {
	return time.AsTime()
}
