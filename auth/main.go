package main

import (
	"auth/db"
	"auth/models"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"

	proto "forum/protos/go"
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
	proto.RegisterAuthServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		os.Exit(1)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// server implements the AuthServiceServer
type server struct {
	proto.UnimplementedAuthServiceServer
}

func (s *server) GetUsernameByUserID(req *proto.UserIDRequest, stream proto.AuthService_GetUsernameByUserIDServer) error {
	username, err := models.GetUsernameByUserID(int(req.GetUserId()))
	if err != nil {
		return err
	}
	return stream.Send(&proto.UsernameResponse{Username: username})
}
