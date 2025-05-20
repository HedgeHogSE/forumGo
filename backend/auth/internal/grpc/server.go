package grpc

import (
	"context"
	"forum/backend/auth/internal/models"
	"forum/backend/protos/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func StartGRPCServer() {
	grpcServer := grpc.NewServer()
	userpb.RegisterAuthServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		os.Exit(1)
	}

	log.Println("gRPC server started on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

type server struct {
	userpb.UnimplementedAuthServiceServer
}

func (s *server) GetUserName(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	username, err := models.GetUsernameByUserID(int(req.GetUserId()))
	if err != nil {
		return nil, err
	}
	return &userpb.UserResponse{UserName: username}, nil
}
