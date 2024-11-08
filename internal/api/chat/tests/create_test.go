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
)

func TestCreate(t *testing.T) {
	type (
		chatServiceMockFunc func(mc *minimock.Controller) servinterfaces.ChatService
		args                struct {
			ctx context.Context
			req *chatv1.CreateRequest
		}
	)

	var (
		usernames   = []string{"Ivan, Vitya, Petya"}
		IDs         = models.IDs{1, 2, 3}
		id          = gofakeit.Int64()
		ctx         = context.Background()
		mc          = minimock.NewController(t)
		errorGetIDs = errors.New("error getting IDs")
		errorCreate = errors.New("error creating chat")
		res         = &chatv1.CreateResponse{
			Id: id,
		}
	)

	defer t.Cleanup(mc.Finish)
	tests := []struct {
		name            string
		args            args
		want            *chatv1.CreateResponse
		err             error
		chatServiceMock chatServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &chatv1.CreateRequest{
					Usernames: usernames,
				},
			},
			want: res,
			err:  nil,
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(IDs, nil)
				mock.CreateMock.Expect(ctx, IDs).Return(&id, nil)
				return mock
			},
		},
		{
			name: "error case: error with GetIDs",
			args: args{
				ctx: ctx,
				req: &chatv1.CreateRequest{
					Usernames: usernames,
				},
			},
			want: nil,
			err:  fmt.Errorf("error with get IDs: %w", errorGetIDs),
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(nil, errorGetIDs)
				return mock
			},
		},
		{
			name: "error case: error with Create",
			args: args{
				ctx: ctx,
				req: &chatv1.CreateRequest{
					Usernames: usernames,
				},
			},
			want: nil,
			err:  fmt.Errorf("error when trying to create a chat: %w", errorCreate),
			chatServiceMock: func(mc *minimock.Controller) servinterfaces.ChatService {
				mock := serviceMocks.NewChatServiceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(IDs, nil)
				mock.CreateMock.Expect(ctx, IDs).Return(nil, errorCreate)
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

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, newID)
		})
	}
}
