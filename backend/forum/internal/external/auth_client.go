package external

import (
	"context"
	"forum/backend/protos/go"
	"google.golang.org/grpc"
	"log"
)

func GetUsernameFromAuth(userID int) (string, error) {
	// Устанавливаем соединение с gRPC-сервером auth
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to auth service: %v", err)
		return "", err
	}
	defer conn.Close()

	// Создаём клиента
	client := userpb.NewAuthServiceClient(conn)

	// Отправляем запрос
	resp, err := client.GetUserName(context.Background(), &userpb.UserRequest{UserId: int32(userID)})
	if err != nil {
		log.Printf("Error calling GetUsernameByUserID: %v", err)
		return "", err
	}

	return resp.GetUserName(), nil
}
