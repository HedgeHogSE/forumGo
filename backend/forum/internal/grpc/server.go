package grpc

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/HedgeHogSE/forum/backend/forum/internal/models"
	userpb "github.com/HedgeHogSE/forum/backend/protos/go"

	"google.golang.org/grpc"
)

// CommentService определяет интерфейс для работы с комментариями
type CommentService interface {
	GetCommentsByAuthorID(authorID int) ([]models.Comment, error)
}

type BackendServer struct {
	userpb.UnimplementedBackendServiceServer
	commentService CommentService
}

// NewBackendServer создает новый экземпляр BackendServer
func NewBackendServer(commentService CommentService) *BackendServer {
	return &BackendServer{
		commentService: commentService,
	}
}

func (s *BackendServer) GetUserComments(ctx context.Context, req *userpb.UserCommentsRequest) (*userpb.UserCommentsResponse, error) {
	comments, err := s.commentService.GetCommentsByAuthorID(int(req.UserId))
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
	userpb.RegisterBackendServiceServer(s, NewBackendServer(models.NewCommentService()))

	log.Printf("Backend gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
