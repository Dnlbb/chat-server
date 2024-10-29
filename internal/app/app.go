package app

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Dnlbb/chat-server/internal/closer"
	"github.com/Dnlbb/chat-server/internal/config"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App структура сервисной модели
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

// NewApp конструктор для сервисной модели
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, fmt.Errorf("init deps: %w", err)
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if err := config.LoadEnv("chat.env"); err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	chatv1.RegisterChatServer(a.grpcServer, a.serviceProvider.GetChatController(ctx))

	return nil
}

// Run старт
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer()
}

func (a *App) runGRPCServer() error {
	log.Printf("starting gRPC server on %s", a.serviceProvider.GetGRPCConfig().Address())

	listener, err := net.Listen("tcp", a.serviceProvider.GetGRPCConfig().Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
