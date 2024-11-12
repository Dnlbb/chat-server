package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	clientMocks "github.com/Dnlbb/chat-server/internal/client/mocks"
	"github.com/Dnlbb/chat-server/internal/models"
	repoMocks "github.com/Dnlbb/chat-server/internal/repository/mocks"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/chatserv"
	"github.com/Dnlbb/platform_common/pkg/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()

	type (
		ChatTxManMockFunc   func(mc *minimock.Controller) db.TxManager
		ChatStorageMockFunc func(mc *minimock.Controller) repointerface.StorageInterface
		AuthStorageMockFunc func(mc *minimock.Controller) repointerface.AuthInterface
		args                struct {
			ctx     context.Context
			message models.Message
		}
	)

	var (
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		errLog    = errors.New("log error")
		errSend   = errors.New("send error")
		chatID    = gofakeit.Int64()
		fromUName = gofakeit.Name()
		fromUID   = gofakeit.Int64()
		Body      = gofakeit.PetName()
		time      = time.Now()
		message   = models.Message{
			ChatID:    chatID,
			FromUname: fromUName,
			FromUID:   fromUID,
			Body:      Body,
			Time:      time,
		}
	)

	defer mc.Finish()

	tests := []struct {
		name            string
		args            args
		want            *int64
		err             error
		chatTxManMock   ChatTxManMockFunc
		chatStorageMock ChatStorageMockFunc
		authStorageMock AuthStorageMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: nil,
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.SendMessageChatMock.Expect(ctx, message).Return(nil)
				mock.LogMock.Expect(ctx, models.SENDMESSAGE).Return(nil)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				return repoMocks.NewAuthInterfaceMock(mc)
			},
		},
		{
			name: "error case: error with send message from storage",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: fmt.Errorf("error while send message chat: %w", errSend),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.SendMessageChatMock.Expect(ctx, message).Return(errSend)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				return repoMocks.NewAuthInterfaceMock(mc)
			},
		},
		{
			name: "error case: error with logging",
			args: args{
				ctx:     ctx,
				message: message,
			},
			err: fmt.Errorf("error log: %w", errLog),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.SendMessageChatMock.Expect(ctx, message).Return(nil)
				mock.LogMock.Expect(ctx, models.SENDMESSAGE).Return(errLog)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				return repoMocks.NewAuthInterfaceMock(mc)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			RepoMock := tt.chatStorageMock(mc)
			TxManMock := tt.chatTxManMock(mc)
			authClient := tt.authStorageMock(mc)
			service := chatserv.NewService(RepoMock, TxManMock, authClient)

			err := service.SendMessage(tt.args.ctx, tt.args.message)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
