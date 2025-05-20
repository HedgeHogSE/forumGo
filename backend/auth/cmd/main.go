package main

import (
	"context"
	"forum/backend/auth/db"
	"forum/backend/auth/models"
	"forum/backend/protos/go"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

var router *gin.Engine

func main() {
	db.SetupDB()

	go startGRPCServer()

	router = gin.Default()

	initializeRoutes()

	router.Run(":8081")
}

func startGRPCServer() {
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
