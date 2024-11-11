package chat

import (
	"time"

	"github.com/Dnlbb/chat-server/internal/models"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// toChatID конвертация в бизнес модель Chat.
func toChatID(req *chatv1.DeleteRequest) models.Chat {
	return models.Chat{
		ID: req.Id,
	}
}

// toModelsMessage конвертация в бизнес модель Message.
func toMessage(req *chatv1.SendMessageRequest) *models.Message {
	return &models.Message{
		ChatID:    req.GetChatID(),
		FromUname: req.GetFromUserName(),
		FromUID:   req.GetFromUserID(),
		Body:      req.GetBody(),
		Time:      toTimestampTime(req.GetTime()),
	}
}

// ToTimestampProto конвертация proto времени в обычный time.Time.
func toTimestampTime(time *timestamppb.Timestamp) time.Time {
	return time.AsTime()
}

func toUsernames(req *chatv1.CreateRequest) models.Usernames {
	usernames := append(models.Usernames{}, req.Usernames...)
	return usernames
}
