package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Dnlbb/chat-server/internal/api/chat"
	"github.com/Dnlbb/chat-server/internal/models"
	serviceMocks "github.com/Dnlbb/chat-server/internal/service/mocks"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestSendMessage(t *testing.T) {
	type (
		chatServiceMockFunc func(mc *minimock.Controller) servinterfaces.ChatService
		args                struct {
			ctx context.Context
			req *chatv1.SendMessageRequest
		}
	)

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		ID      = gofakeit.Int64()
		uname   = gofakeit.Username()
		UID     = gofakeit.Int64()
		body    = gofakeit.City()
		time    = gofakeit.Date()
		message = models.Message{
			ChatID:    ID,
			FromUname: uname,
			FromUID:   UID,
			Body:      body,
			Time:      time,
		}
		errSend = errors.New("send error")
	)

	defer t.Cleanup(mc.Finish)
	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &chatv1.SendMessageRequest{
					ChatID:    ID,
					FromUname: uname,
					FromUid:   UID,
					Body:      body,
					Time:      timestamppb.New(time),
				},
			},
			want: &emptypb.Empty{},
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(nil)
				return mock
			},
		},
		{
			name: "error case: error with send message from chat service",
			args: args{
				ctx: ctx,
				req: &chatv1.SendMessageRequest{
					ChatID:    ID,
					FromUname: uname,
					FromUid:   UID,
					Body:      body,
					Time:      timestamppb.New(time),
				},
			},
			want: &emptypb.Empty{},
			err:  fmt.Errorf("failed to send message: %w", errSend),
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.SendMessageMock.Expect(ctx, message).Return(errSend)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ChatServiceMock := tt.chatServiceMock(mc)
			api := chat.NewController(ChatServiceMock)

			newID, err := api.SendMessage(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, newID)
		})
	}
}
