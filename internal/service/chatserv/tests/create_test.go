package tests

import (
	"context"
	"fmt"
	"testing"

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

func TestCreate(t *testing.T) {
	t.Parallel()

	type (
		ChatTxManMockFunc   func(mc *minimock.Controller) db.TxManager
		ChatStorageMockFunc func(mc *minimock.Controller) repointerface.StorageInterface
		AuthStorageMockFunc func(mc *minimock.Controller) repointerface.AuthInterface
		args                struct {
			ctx       context.Context
			usernames models.Usernames
		}
	)

	var (
		ctx           = context.Background()
		mc            = minimock.NewController(t)
		usernames     = models.Usernames{"Ivan", "Petr", "Viktor"}
		IDs           = models.IDs{1, 2, 3}
		ChatID        = gofakeit.Int64()
		errCreateChat = errors.New("create chat error")
		errGetIDs     = errors.New("get IDs error")
		errLog        = errors.New("log error")
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
				ctx:       ctx,
				usernames: usernames,
			},
			want: &ChatID,
			err:  nil,
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.CreateChatMock.Expect(ctx, IDs).Return(&ChatID, nil)
				mock.LogMock.Expect(ctx, models.CREATE).Return(nil)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				mock := repoMocks.NewAuthInterfaceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(IDs, nil)
				return mock
			},
		},
		{
			name: "error case: error with GetIDs",
			args: args{
				ctx:       ctx,
				usernames: usernames,
			},
			want: nil,
			err:  fmt.Errorf("error with get IDs: %w", errGetIDs),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				mock := repoMocks.NewAuthInterfaceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(nil, errGetIDs)
				return mock
			},
		},
		{
			name: "error case: error with create chat",
			args: args{
				ctx:       ctx,
				usernames: usernames,
			},
			want: nil,
			err:  fmt.Errorf("error creating chat: %w", errCreateChat),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.CreateChatMock.Expect(ctx, IDs).Return(nil, errCreateChat)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				mock := repoMocks.NewAuthInterfaceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(IDs, nil)
				return mock
			},
		},
		{
			name: "error case: error with logging create",
			args: args{
				ctx:       ctx,
				usernames: usernames,
			},
			want: nil,
			err:  fmt.Errorf("error logging create chat: %w", errLog),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				mock.ReadCommittedMock.Set(func(ctx context.Context, handler db.Handler) error {
					return handler(ctx)
				})
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				mock.CreateChatMock.Expect(ctx, IDs).Return(&ChatID, nil)
				mock.LogMock.Expect(ctx, models.CREATE).Return(errLog)
				return mock
			},
			authStorageMock: func(mc *minimock.Controller) repointerface.AuthInterface {
				mock := repoMocks.NewAuthInterfaceMock(mc)
				mock.GetIDsMock.Expect(ctx, usernames).Return(IDs, nil)
				return mock
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

			id, err := service.Create(tt.args.ctx, tt.args.usernames)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, id)
		})
	}
}
