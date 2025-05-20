package external

import (
	"context"
	"forum/backend/protos/go"
	"log"

	"google.golang.org/grpc"
)

func GetUserCommentsFromBackend(userID int) ([]*userpb.Comment, error) {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to forum service: %v", err)
		return nil, err
	}
	defer conn.Close()

	client := userpb.NewBackendServiceClient(conn)

	resp, err := client.GetUserComments(context.Background(), &userpb.UserCommentsRequest{UserId: int32(userID)})
	if err != nil {
		log.Printf("Error calling GetUserComments: %v", err)
		return nil, err
	}

	return resp.GetComments(), nil
}
