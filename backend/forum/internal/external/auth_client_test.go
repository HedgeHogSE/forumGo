package external_test

import (
	"context"
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

func setupTestServer(t *testing.T, mockServer *MockAuthServer) (*grpc.Server, func()) {
	// Создаем тестовый сервер
	lis, err := net.Listen("tcp", ":50051") // Используем тот же порт, что и в оригинальной функции
	require.NoError(t, err)

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

	return s, cleanup
}

func TestGetUsernameByUserID(t *testing.T) {
	// Создаем мок сервера
	mockServer := &MockAuthServer{
		username: "testuser",
	}

	// Настраиваем тестовый сервер
	_, cleanup := setupTestServer(t, mockServer)
	defer cleanup()

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
	_, cleanup := setupTestServer(t, mockServer)
	defer cleanup()

	// Тестируем получение имени пользователя с ошибкой
	_, err := external.GetUsernameByUserID(1)
	assert.Error(t, err)
}

func TestGetUsernameByUserID_ConnectionError(t *testing.T) {
	// Тестируем ошибку подключения к серверу
	_, err := external.GetUsernameByUserID(1)
	assert.Error(t, err)
}
