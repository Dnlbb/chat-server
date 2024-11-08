package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Dnlbb/auth/pkg/auth_v1"
	authv1 "github.com/Dnlbb/auth/pkg/auth_v1"
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
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetIDs(t *testing.T) {
	t.Parallel()

	type (
		ChatTxManMockFunc   func(mc *minimock.Controller) db.TxManager
		ChatStorageMockFunc func(mc *minimock.Controller) repointerface.StorageInterface
		AuthClientMockFunc  func(mc *minimock.Controller) auth_v1.AuthClient
		args                struct {
			ctx       context.Context
			Usernames models.Usernames
		}
	)

	var (
		ctx         = context.Background()
		mc          = minimock.NewController(t)
		Usernames   = models.Usernames{"Petr", "Ivan", "Viktor"}
		IDs         = models.IDs{1, 2, 3}
		email       = gofakeit.Email()
		role        = authv1.Role_USER
		createdAt   = gofakeit.Date()
		updatedAt   = gofakeit.Date()
		errFromAuth = errors.New("error from Auth")
	)

	defer mc.Finish()

	tests := []struct {
		name            string
		args            args
		want            models.IDs
		err             error
		chatTxManMock   ChatTxManMockFunc
		chatStorageMock ChatStorageMockFunc
		authClientMock  AuthClientMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:       ctx,
				Usernames: Usernames,
			},
			want: IDs,
			err:  nil,
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				return mock
			},
			authClientMock: func(mc *minimock.Controller) auth_v1.AuthClient {
				mock := mocks.NewAuthClientMock(mc)
				for i, username := range Usernames {
					id := int64(i + 1)
					mock.GetMock.When(ctx, &authv1.GetRequest{NameOrId: &authv1.GetRequest_Username{
						Username: username,
					}}).Then(&authv1.GetResponse{Id: id,
						User: &authv1.User{
							Name:  username,
							Email: email,
							Role:  role,
						},
						CreatedAt: timestamppb.New(createdAt),
						UpdatedAt: timestamppb.New(updatedAt),
					}, nil)
				}
				return mock
			},
		},
		{
			name: "error case: error with client auth",
			args: args{
				ctx:       ctx,
				Usernames: Usernames,
			},
			want: nil,
			err:  fmt.Errorf("error when trying to get a user profile from the authorization service: %w", errFromAuth),
			chatTxManMock: func(mc *minimock.Controller) db.TxManager {
				mock := clientMocks.NewTxManagerMock(mc)
				return mock
			},
			chatStorageMock: func(mc *minimock.Controller) repointerface.StorageInterface {
				mock := repoMocks.NewStorageInterfaceMock(mc)
				return mock
			},
			authClientMock: func(mc *minimock.Controller) auth_v1.AuthClient {
				mock := mocks.NewAuthClientMock(mc)
				mock.GetMock.When(ctx, &authv1.GetRequest{NameOrId: &authv1.GetRequest_Username{
					Username: Usernames[0],
				}}).Then(nil, errFromAuth)
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
			authClient := tt.authClientMock(mc)
			service := chatserv.NewService(RepoMock, TxManMock, authClient)

			IDs, err := service.GetIDs(tt.args.ctx, tt.args.Usernames)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, IDs)
		})
	}
}
