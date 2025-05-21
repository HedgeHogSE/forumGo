package external

import (
	"context"
	"log"

	userpb "github.com/HedgeHogSE/forum/backend/protos/go"

	"google.golang.org/grpc"
)

// GetUsernameByUserIDFunc - тип функции для получения имени пользователя
type GetUsernameByUserIDFunc func(userID int) (string, error)

// GetUsernameByUserID - функция для получения имени пользователя по ID
var GetUsernameByUserID GetUsernameByUserIDFunc = func(userID int) (string, error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to auth service: %v", err)
		return "", err
	}
	defer conn.Close()

	client := userpb.NewAuthServiceClient(conn)

	resp, err := client.GetUserName(context.Background(), &userpb.UserRequest{UserId: int32(userID)})
	if err != nil {
		log.Printf("Error calling GetUserName: %v", err)
		return "", err
	}

	return resp.GetUserName(), nil
}
