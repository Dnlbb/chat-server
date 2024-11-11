package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	chatapi "github.com/Dnlbb/chat-server/internal/api/chat"
	"github.com/Dnlbb/chat-server/internal/models"
	serviceMocks "github.com/Dnlbb/chat-server/internal/service/mocks"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelete(t *testing.T) {
	type (
		chatServiceMockFunc func(mc *minimock.Controller) servinterfaces.ChatService
		args                struct {
			ctx context.Context
			req *chatv1.DeleteRequest
		}
	)

	var (
		ID   = gofakeit.Int64()
		chat = models.Chat{
			ID: ID,
		}
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		errDelete = errors.New("error with delete")
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
				req: &chatv1.DeleteRequest{
					Id: ID,
				},
			},
			want: &emptypb.Empty{},
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, chat).Return(nil)
				return mock
			},
		},
		{
			name: "error case: error with delete from chat service",
			args: args{
				ctx: ctx,
				req: &chatv1.DeleteRequest{
					Id: ID,
				},
			},
			want: &emptypb.Empty{},
			err:  fmt.Errorf("delete chat error: %w", errDelete),
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.DeleteMock.Expect(ctx, chat).Return(errDelete)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ChatServiceMock := tt.chatServiceMock(mc)
			api := chatapi.NewController(ChatServiceMock)

			newID, err := api.Delete(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, newID)
		})
	}
}
