package grpc

import (
	"context"
	"forum/backend/forum/internal/models"
	"forum/backend/protos/go"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type BackendServer struct {
	userpb.UnimplementedBackendServiceServer
}

func (s *BackendServer) GetUserComments(ctx context.Context, req *userpb.UserCommentsRequest) (*userpb.UserCommentsResponse, error) {
	comments, err := models.GetCommentsByAuthorID(int(req.UserId))
	if err != nil {
		return nil, err
	}

	protoComments := make([]*userpb.Comment, len(comments))
	for i, comment := range comments {
		protoComments[i] = &userpb.Comment{
			Id:        int32(comment.ID),
			Content:   comment.Content,
			TopicId:   int32(comment.TopicId),
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		}
	}

	return &userpb.UserCommentsResponse{
		Comments: protoComments,
	}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userpb.RegisterBackendServiceServer(s, &BackendServer{})

	log.Printf("Backend gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
