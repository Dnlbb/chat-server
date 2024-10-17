package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Dnlbb/chat-server/postgres/postgresstorage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	authv1 "github.com/Dnlbb/auth/pkg/auth_v1"

	desc "github.com/Dnlbb/chat-server/pkg/chat"
)

func main() {
	authConn, err := grpc.Dial("127.0.0.1:50501", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Не удалось установить соединение с сервисом авторизации: %v", err)
	}
	defer closer(authConn)
	authClient := authv1.NewAuthClient(authConn)

	lis, err := net.Listen("tcp", "127.0.0.1:50050")
	if err != nil {
		log.Fatal("Ошибка при старте сервера на порту 50050")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterShatServer(s, &server{authClient: authClient})

	log.Printf("Cервер запущен на  %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	desc.UnimplementedShatServer
	authClient authv1.AuthClient
	storage    postgresstorage.StorageInterface
}

func (s *server) CreateChat(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Users: %v", req.GetUsernames())
	IDs := make(postgresstorage.IDs, len(req.GetUsernames()))
	for _, username := range req.GetUsernames() {
		resp, err := s.authClient.Get(context.Background(), &authv1.GetRequest{NameOrId: &authv1.GetRequest_Username{
			Username: username,
		}})
		if err != nil {
			return nil, fmt.Errorf("error when trying to get a user profile from the authorization service: %w", err)
		}
		IDs = append(IDs, resp.Id)
		log.Printf("Получен профиль пользователя %s: %v", username, resp)
	}
	id, err := s.storage.CreateChat(postgresstorage.IDs(IDs))
	if err != nil {
		return nil, fmt.Errorf("error when trying to create a chat: %w", err)
	}

	return &desc.CreateResponse{Id: *id}, nil
}

func (s *server) DeleteChat(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Id: %d", req.GetId())
	return nil, nil
}

func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("From: %s", req.GetFrom())
	log.Printf("Text: %s", req.GetText())
	log.Printf("Time: %v", req.GetTimestamp())
	return nil, nil
}

func closer(authConn *grpc.ClientConn) {
	err := authConn.Close()
	if err != nil {
		log.Fatalf("error when closing connection")
	}
}
