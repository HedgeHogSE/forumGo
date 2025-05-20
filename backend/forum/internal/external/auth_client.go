package external

import (
	"context"
	userpb "forum/backend/protos/go"
	"log"

	"google.golang.org/grpc"
)

func GetUsernameFromAuth(userID int) (string, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to auth service: %v", err)
		return "", err
	}
	defer conn.Close()

	client := userpb.NewAuthServiceClient(conn)

	resp, err := client.GetUserName(context.Background(), &userpb.UserRequest{UserId: int32(userID)})
	if err != nil {
		log.Printf("Error calling GetUsernameByUserID: %v", err)
		return "", err
	}

	return resp.GetUserName(), nil
}
