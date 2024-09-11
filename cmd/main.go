package main

import (
	"context"
	"log"
	"net"

	desc "github.com/Dnlbb/chat-server/pkg/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	desc.UnimplementedShatServer
}

func (s *server) Create(_ context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("Users: %v", req.GetUsernames())
	return &desc.CreateResponse{}, nil
}

func (s *server) Delete(_ context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Id: %d", req.GetId())
	return nil, nil
}

func (s *server) SendMessage(_ context.Context, req *desc.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("From: %s", req.GetFrom())
	log.Printf("Text: %s", req.GetText())
	log.Printf("Time: %v", req.GetTimestamp())
	return nil, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:50050")
	if err != nil {
		log.Fatal("failed to listen: 50050 ")
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterShatServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
