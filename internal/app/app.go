package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Dnlbb/chat-server/internal/interceptor"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/Dnlbb/chat-server/internal/config"
	chatv1 "github.com/Dnlbb/chat-server/pkg/chat_v1"
	_ "github.com/Dnlbb/chat-server/statik" // Нужно для инициализации файловой системы.
	"github.com/Dnlbb/platform_common/pkg/closer"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App структура сервисной модели
type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

// NewApp конструктор для приложения модели
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
		a.initHTTPServer,
		a.initSwaggerServer,
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
	interceptors := grpc.ChainUnaryInterceptor(interceptor.ValidateInterceptor, a.serviceProvider.GetAuthInterceptor(ctx).AccessInterceptor)

	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()), interceptors)

	reflection.Register(a.grpcServer)

	chatv1.RegisterChatServer(a.grpcServer, a.serviceProvider.GetChatController(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := chatv1.RegisterChatHandlerFromEndpoint(ctx, mux, a.serviceProvider.GetGRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Content-Length", "Accept"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.GetHTTPConfig().Address(),
		Handler:           cors.Handler(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("init statikfs: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	mux.HandleFunc("/api.swagger.json", SwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.GetSwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return nil
}

// SwaggerFile реализация документации.
func SwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		statikFS, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		content, err := statikFS.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		closer.Add(content.Close)

		if _, err := io.Copy(w, content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

// Run старт
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
		if err != nil {
			log.Printf("grpc server error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.runHTTPServer()
		if err != nil {
			log.Printf("http server error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.runSwaggerServer()
		if err != nil {
			log.Printf("swagger server error: %v", err)
		}
	}()

	wg.Wait()

	return nil
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

func (a *App) runHTTPServer() error {
	log.Printf("starting HTTP server on %s", a.serviceProvider.GetHTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.GetSwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
