package app

import (
	"context"
	"log"

	authv1 "github.com/Dnlbb/auth/pkg/auth_v1"
	"github.com/Dnlbb/chat-server/internal/api/chat"
	"github.com/Dnlbb/chat-server/internal/client/db"
	"github.com/Dnlbb/chat-server/internal/client/db/pg"
	"github.com/Dnlbb/chat-server/internal/client/db/transaction"
	"github.com/Dnlbb/chat-server/internal/closer"
	"github.com/Dnlbb/chat-server/internal/config"
	"github.com/Dnlbb/chat-server/internal/repository/postgres/storage"
	"github.com/Dnlbb/chat-server/internal/repository/repointerface"
	"github.com/Dnlbb/chat-server/internal/service/chatserv"
	"github.com/Dnlbb/chat-server/internal/service/servinterfaces"
	"google.golang.org/grpc"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	txManager      db.TxManager
	chatRepository repointerface.StorageInterface

	chatService servinterfaces.ChatService
	authClient  authv1.AuthClient

	authController *chat.Controller
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// GetPGConfig получаем конфиг постгреса, для последующего получения DSN.
func (s *serviceProvider) GetPGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPgConfig()
		if err != nil {
			log.Fatal("failed to load pg config: %w", err)
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

// GetGRPCConfig получаем конфиг grpc, для последующего получения адреса нашего приложения и старта на нем.
func (s *serviceProvider) GetGRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGrpcConfig()
		if err != nil {
			log.Fatal("failed to load gRPC config: %w", err)
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

// GetDBClient инициализируем клиента к базе данных.
func (s *serviceProvider) GetDBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.GetPGConfig().DSN())
		if err != nil {
			log.Fatal("failed to connect to database: %w", err)
		}

		if err = cl.DB().Ping(ctx); err != nil {
			log.Fatal("failed to ping database: %w", err)
		}

		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

// GetTxManager инициализация менеджера транзакций.
func (s *serviceProvider) GetTxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.GetDBClient(ctx).DB())
	}

	return s.txManager
}

// GetAuthRepository инициализация хранилища.
func (s *serviceProvider) GetAuthRepository(ctx context.Context) repointerface.StorageInterface {
	if s.chatRepository == nil {
		s.chatRepository = storage.NewPostgresRepo(s.GetDBClient(ctx))
	}

	return s.chatRepository
}

// GetAuthService инициализация сервиса авторизации.
func (s *serviceProvider) GetAuthService(ctx context.Context) servinterfaces.ChatService {
	if s.chatService == nil {
		s.chatService = chatserv.NewService(s.GetAuthRepository(ctx),
			s.GetTxManager(ctx),
		)
	}

	return s.chatService
}

func (s *serviceProvider) GetAuthClient() authv1.AuthClient {
	if s.authClient == nil {
		conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		closer.Add(conn.Close)

		s.authClient = authv1.NewAuthClient(conn)
	}

	return s.authClient
}

// GetChatController инициализация контроллера.
func (s *serviceProvider) GetChatController(ctx context.Context) *chat.Controller {
	if s.authController == nil {

		s.authController = chat.NewController(s.GetAuthService(ctx), s.GetAuthClient())
	}

	return s.authController
}