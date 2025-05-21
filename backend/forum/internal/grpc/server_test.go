package grpc_test

import (
	"context"
	"net"
	"os"
	"testing"
	"time"

	"github.com/HedgeHogSE/forum/backend/forum/internal/db"
	grpcserver "github.com/HedgeHogSE/forum/backend/forum/internal/grpc"
	"github.com/HedgeHogSE/forum/backend/forum/internal/models"
	userpb "github.com/HedgeHogSE/forum/backend/protos/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func setupTestDB(t *testing.T) {
	// Устанавливаем тестовые значения для БД
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "Sashaezhak2006")
	os.Setenv("DB_NAME", "forum_test")
	os.Setenv("DB_SSLMODE", "disable")

	// Подключаемся к тестовой БД
	err := db.SetupDB()
	require.NoError(t, err)

	// Очищаем таблицы перед тестами
	_, err = db.Db.Exec("TRUNCATE comments, topics, users CASCADE")
	require.NoError(t, err)
}

// MockCommentService - мок для тестирования
type MockCommentService struct {
	comments []models.Comment
	err      error
}

func (m *MockCommentService) GetCommentsByAuthorID(authorID int) ([]models.Comment, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.comments, nil
}

func setupTestServer(t *testing.T, mockService *MockCommentService) (*grpc.Server, *userpb.BackendServiceClient, func()) {
	// Создаем тестовый сервер
	lis, err := net.Listen("tcp", ":0") // Используем случайный порт
	require.NoError(t, err)

	s := grpc.NewServer()
	userpb.RegisterBackendServiceServer(s, grpcserver.NewBackendServer(mockService))

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("failed to serve: %v", err)
		}
	}()

	// Подключаемся к серверу
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)

	client := userpb.NewBackendServiceClient(conn)

	// Возвращаем функцию очистки
	cleanup := func() {
		s.Stop()
		conn.Close()
	}

	return s, &client, cleanup
}

func TestGetUserComments(t *testing.T) {
	// Создаем тестовые комментарии
	testComments := []models.Comment{
		{
			ID:        1,
			Content:   "Test Comment 1",
			TopicId:   1,
			CreatedAt: time.Now(),
		},
		{
			ID:        2,
			Content:   "Test Comment 2",
			TopicId:   1,
			CreatedAt: time.Now(),
		},
	}

	// Создаем мок сервиса
	mockService := &MockCommentService{
		comments: testComments,
	}

	// Настраиваем тестовый сервер
	_, client, cleanup := setupTestServer(t, mockService)
	defer cleanup()

	// Тестируем получение комментариев
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := (*client).GetUserComments(ctx, &userpb.UserCommentsRequest{
		UserId: 1,
	})
	require.NoError(t, err)

	// Проверяем результаты
	assert.Len(t, resp.Comments, 2)
	assert.Equal(t, "Test Comment 1", resp.Comments[0].Content)
	assert.Equal(t, "Test Comment 2", resp.Comments[1].Content)
	assert.Equal(t, int32(1), resp.Comments[0].TopicId)
	assert.Equal(t, int32(1), resp.Comments[1].TopicId)
}

func TestGetUserComments_NoComments(t *testing.T) {
	// Создаем мок сервиса с пустым списком комментариев
	mockService := &MockCommentService{
		comments: []models.Comment{},
	}

	// Настраиваем тестовый сервер
	_, client, cleanup := setupTestServer(t, mockService)
	defer cleanup()

	// Тестируем получение комментариев для пользователя без комментариев
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := (*client).GetUserComments(ctx, &userpb.UserCommentsRequest{
		UserId: 1,
	})
	require.NoError(t, err)

	// Проверяем результаты
	assert.Empty(t, resp.Comments)
}

func TestGetUserComments_Error(t *testing.T) {
	// Создаем мок сервиса с ошибкой
	mockService := &MockCommentService{
		err: assert.AnError,
	}

	// Настраиваем тестовый сервер
	_, client, cleanup := setupTestServer(t, mockService)
	defer cleanup()

	// Тестируем получение комментариев с ошибкой
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := (*client).GetUserComments(ctx, &userpb.UserCommentsRequest{
		UserId: 1,
	})
	assert.Error(t, err)
}

func TestGetUserComments_InvalidUser(t *testing.T) {
	// Создаем мок сервиса с пустым списком комментариев
	mockService := &MockCommentService{
		comments: []models.Comment{},
	}

	// Настраиваем тестовый сервер
	_, client, cleanup := setupTestServer(t, mockService)
	defer cleanup()

	// Тестируем получение комментариев для несуществующего пользователя
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := (*client).GetUserComments(ctx, &userpb.UserCommentsRequest{
		UserId: 999, // Несуществующий ID
	})
	require.NoError(t, err)

	// Проверяем результаты
	assert.Empty(t, resp.Comments)
}
