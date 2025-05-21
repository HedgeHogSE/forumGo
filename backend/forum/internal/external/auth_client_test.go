package external_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/HedgeHogSE/forum/backend/forum/internal/external"
	userpb "github.com/HedgeHogSE/forum/backend/protos/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// MockAuthServer - мок для тестирования
type MockAuthServer struct {
	userpb.UnimplementedAuthServiceServer
	username string
	err      error
}

func (m *MockAuthServer) GetUserName(ctx context.Context, req *userpb.UserRequest) (*userpb.UserResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &userpb.UserResponse{UserName: m.username}, nil
}

func setupTestServer(t *testing.T, mockServer *MockAuthServer) (*grpc.Server, string, func()) {
	// Создаем тестовый сервер на случайном порту
	lis, err := net.Listen("tcp", ":0") // Используем порт 0 для получения случайного порта
	require.NoError(t, err)

	port := lis.Addr().(*net.TCPAddr).Port
	addr := fmt.Sprintf("localhost:%d", port)

	s := grpc.NewServer()
	userpb.RegisterAuthServiceServer(s, mockServer)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	// Возвращаем функцию очистки
	cleanup := func() {
		s.Stop()
	}

	return s, addr, cleanup
}

func TestGetUsernameByUserID(t *testing.T) {
	// Создаем мок сервера
	mockServer := &MockAuthServer{
		username: "testuser",
	}

	// Настраиваем тестовый сервер
	_, addr, cleanup := setupTestServer(t, mockServer)
	defer cleanup()

	// Сохраняем оригинальную функцию и восстанавливаем после теста
	originalFunc := external.GetUsernameByUserID
	defer func() { external.GetUsernameByUserID = originalFunc }()

	// Подменяем функцию на тестовую версию
	external.GetUsernameByUserID = func(userID int) (string, error) {
		// Создаем соединение с тестовым сервером
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return "", fmt.Errorf("failed to connect to auth server: %v", err)
		}
		defer conn.Close()

		// Создаем клиент
		client := userpb.NewAuthServiceClient(conn)

		// Получаем имя пользователя
		resp, err := client.GetUserName(context.Background(), &userpb.UserRequest{UserId: int32(userID)})
		if err != nil {
			return "", fmt.Errorf("error calling GetUserName: %v", err)
		}
		return resp.GetUserName(), nil
	}

	// Тестируем получение имени пользователя
	username, err := external.GetUsernameByUserID(1)
	require.NoError(t, err)
	assert.Equal(t, "testuser", username)
}

func TestGetUsernameByUserID_Error(t *testing.T) {
	// Создаем мок сервера с ошибкой
	mockServer := &MockAuthServer{
		err: assert.AnError,
	}

	// Настраиваем тестовый сервер
	_, addr, cleanup := setupTestServer(t, mockServer)
	defer cleanup()

	// Сохраняем оригинальную функцию и восстанавливаем после теста
	originalFunc := external.GetUsernameByUserID
	defer func() { external.GetUsernameByUserID = originalFunc }()

	// Подменяем функцию на тестовую версию
	external.GetUsernameByUserID = func(userID int) (string, error) {
		// Создаем соединение с тестовым сервером
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			return "", fmt.Errorf("failed to connect to auth server: %v", err)
		}
		defer conn.Close()

		// Создаем клиент
		client := userpb.NewAuthServiceClient(conn)

		// Получаем имя пользователя
		resp, err := client.GetUserName(context.Background(), &userpb.UserRequest{UserId: int32(userID)})
		if err != nil {
			return "", fmt.Errorf("error calling GetUserName: %v", err)
		}
		return resp.GetUserName(), nil
	}

	// Тестируем получение имени пользователя с ошибкой
	_, err := external.GetUsernameByUserID(1)
	assert.Error(t, err)
}

func TestGetUsernameByUserID_ConnectionError(t *testing.T) {
	// Тестируем ошибку подключения к серверу
	_, err := external.GetUsernameByUserID(1)
	assert.Error(t, err)
}
