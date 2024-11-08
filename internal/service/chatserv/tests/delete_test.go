package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Dnlbb/auth/pkg/auth_v1"
	clientMocks "github.com/Dnlbb/chat-server/internal/client/mocks"
	"github.com/Dnlbb/chat-server/internal/models"
	repoMocks "github.com/Dnlbb/chat-server/internal/repository/mocks"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/chatserv"
	"github.com/Dnlbb/chat-server/internal/service/mocks"
	"github.com/Dnlbb/platform_common/pkg/db"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type (
		ChatTxManMockFunc   func(mc *minimock.Controller) db.TxManager
		ChatStorageMockFunc func(mc *minimock.Controller) repointerface.StorageInterface
		AuthClientMockFunc  func(mc *minimock.Controller) auth_v1.AuthClient
		args                struct {
			ctx    context.Context
			chatID models.ChatID
		}
	)

	var (
		ctx           = context.Background()
		mc            = minimock.NewController(t)
		chatID        = models.ChatID{ID: gofakeit.Int64()}
		errDeleteChat = errors.New("create delete error")
		errLog        = errors.New("log error")
	)

	defer mc.Finish()

	tests := []struct {
		name            string
		args            args
		err             error
		chatTxManMock   ChatTxManMockFunc
		chatStorageMock ChatStorageMockFunc
		authClientMock  AuthClientMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:    ctx,
				chatID: chatID,
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
				mock.DeleteChatMock.Expect(ctx, chatID).Return(nil)
				mock.LogMock.Expect(ctx, models.DELETE).Return(nil)
				return mock
			},
			authClientMock: func(mc *minimock.Controller) auth_v1.AuthClient {
				return mocks.NewAuthClientMock(mc)
			},
		},
		{
			name: "error case: error with delete chat",
			args: args{
				ctx:    ctx,
				chatID: chatID,
			},
			err: fmt.Errorf("error delete chat: %w", errDeleteChat),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.DeleteChatMock.Expect(ctx, chatID).Return(errDeleteChat)
				return mock
			},
			authClientMock: func(mc *minimock.Controller) auth_v1.AuthClient {
				return mocks.NewAuthClientMock(mc)
			},
		},
		{
			name: "error case: error with logging create",
			args: args{
				ctx:    ctx,
				chatID: chatID,
			},
			err: fmt.Errorf("error logging delete chat: %w", errLog),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.DeleteChatMock.Expect(ctx, chatID).Return(nil)
				mock.LogMock.Expect(ctx, models.DELETE).Return(errLog)
				return mock
			},
			authClientMock: func(mc *minimock.Controller) auth_v1.AuthClient {
				return mocks.NewAuthClientMock(mc)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			RepoMock := tt.chatStorageMock(mc)
			TxManMock := tt.chatTxManMock(mc)
			authClient := tt.authClientMock(mc)
			service := chatserv.NewService(RepoMock, TxManMock, authClient)

			err := service.Delete(tt.args.ctx, tt.args.chatID)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
