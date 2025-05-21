package websocket_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/HedgeHogSE/forum/backend/forum/internal/db"
	"github.com/HedgeHogSE/forum/backend/forum/internal/models"
	ws "github.com/HedgeHogSE/forum/backend/forum/internal/websocket"
	gorilla "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var testUserID int
var testTopicID int

// setupTestServer создает тестовый HTTP сервер с WebSocket handler
func setupTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws.HandleConnections(w, r)
	}))
}

// connectWebSocket подключается к WebSocket серверу с таймаутом
func connectWebSocket(t *testing.T, server *httptest.Server, topicID int) *gorilla.Conn {
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?topic=" + strconv.Itoa(topicID)

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем канал для результата
	result := make(chan *gorilla.Conn)
	errChan := make(chan error)

	go func() {
		ws, _, err := gorilla.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			errChan <- err
			return
		}
		result <- ws
	}()

	// Ждем результат или таймаут
	select {
	case err := <-errChan:
		t.Fatalf("Failed to connect to WebSocket: %v", err)
		return nil
	case ws := <-result:
		return ws
	case <-ctx.Done():
		t.Fatal("Timeout connecting to WebSocket")
		return nil
	}
}

// setupTestDB настраивает тестовую БД
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
	_, err = db.Db.Exec("TRUNCATE comments, topics CASCADE")
	require.NoError(t, err)

	// Создаем тестового пользователя
	createTestUser(t)

	// Создаем тестовый топик
	createTestTopic(t)
}

// createTestUser создает тестового пользователя
func createTestUser(t *testing.T) {
	// Проверяем, существует ли уже тестовый пользователь
	var id int
	err := db.Db.QueryRow("SELECT id FROM users WHERE username = 'test_user'").Scan(&id)
	if err == nil {
		testUserID = id
		return
	}

	// Если пользователь не существует, создаем его
	query := `
		INSERT INTO users (username, email, password_hash, name)
		VALUES ('test_user', 'test@example.com', 'test_hash', 'Test User')
		RETURNING id
	`
	err = db.Db.QueryRow(query).Scan(&testUserID)
	require.NoError(t, err)
}

// createTestTopic создает тестовый топик
func createTestTopic(t *testing.T) {
	query := `
		INSERT INTO topics (title, description, author_id)
		VALUES ('Test Topic', 'Test Description', $1)
		RETURNING id
	`
	err := db.Db.QueryRow(query, testUserID).Scan(&testTopicID)
	require.NoError(t, err)
}

func TestWebSocketConnection(t *testing.T) {
	setupTestDB(t)

	server := setupTestServer(t)
	defer server.Close()

	// Подключаемся к WebSocket
	ws := connectWebSocket(t, server, testTopicID)
	if ws == nil {
		return
	}
	defer ws.Close()

	// Устанавливаем таймаут для чтения
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Проверяем, что получили историю сообщений
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)

	// Проверяем, что сообщение - это JSON массив
	var messages []models.CommentWithUsername
	err = json.Unmarshal(message, &messages)
	require.NoError(t, err)
}

func TestWebSocketMessageExchange(t *testing.T) {
	setupTestDB(t)

	server := setupTestServer(t)
	defer server.Close()

	// Подключаем двух клиентов
	ws1 := connectWebSocket(t, server, testTopicID)
	if ws1 == nil {
		return
	}
	defer ws1.Close()

	ws2 := connectWebSocket(t, server, testTopicID)
	if ws2 == nil {
		return
	}
	defer ws2.Close()

	// Устанавливаем таймауты для чтения
	ws1.SetReadDeadline(time.Now().Add(5 * time.Second))
	ws2.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Пропускаем историю сообщений
	_, _, err := ws1.ReadMessage()
	require.NoError(t, err)
	_, _, err = ws2.ReadMessage()
	require.NoError(t, err)

	// Отправляем сообщение от первого клиента
	message := map[string]interface{}{
		"content":   "test message",
		"topic_id":  testTopicID,
		"author_id": testUserID,
		"username":  "test_user",
	}
	messageBytes, err := json.Marshal(message)
	require.NoError(t, err)

	err = ws1.WriteMessage(gorilla.TextMessage, messageBytes)
	require.NoError(t, err)

	// Проверяем, что второй клиент получил сообщение
	_, receivedMessage, err := ws2.ReadMessage()
	require.NoError(t, err)

	var received map[string]interface{}
	err = json.Unmarshal(receivedMessage, &received)
	require.NoError(t, err)
	assert.Equal(t, message["content"], received["content"])
	assert.Equal(t, float64(message["topic_id"].(int)), received["topic_id"].(float64))
	assert.Equal(t, float64(message["author_id"].(int)), received["author_id"].(float64))
	assert.Equal(t, message["username"], received["username"])
}

func TestWebSocketInvalidTopicID(t *testing.T) {
	setupTestDB(t)

	server := setupTestServer(t)
	defer server.Close()

	// Пытаемся подключиться с неверным topic ID
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?topic=invalid"

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Создаем канал для результата
	errChan := make(chan error)

	go func() {
		_, _, err := gorilla.DefaultDialer.Dial(wsURL, nil)
		errChan <- err
	}()

	// Ждем результат или таймаут
	select {
	case err := <-errChan:
		assert.Error(t, err)
	case <-ctx.Done():
		t.Fatal("Timeout connecting to WebSocket")
	}
}

func TestWebSocketMultipleClients(t *testing.T) {
	setupTestDB(t)

	server := setupTestServer(t)
	defer server.Close()

	// Подключаем трех клиентов
	clients := make([]*gorilla.Conn, 3)
	for i := 0; i < 3; i++ {
		clients[i] = connectWebSocket(t, server, testTopicID)
		if clients[i] == nil {
			return
		}
		defer clients[i].Close()

		// Устанавливаем таймаут для чтения
		clients[i].SetReadDeadline(time.Now().Add(5 * time.Second))

		// Пропускаем историю сообщений
		_, _, err := clients[i].ReadMessage()
		require.NoError(t, err)
	}

	// Отправляем сообщение от первого клиента
	message := map[string]interface{}{
		"content":   "broadcast message",
		"topic_id":  testTopicID,
		"author_id": testUserID,
		"username":  "test_user",
	}
	messageBytes, err := json.Marshal(message)
	require.NoError(t, err)

	err = clients[0].WriteMessage(gorilla.TextMessage, messageBytes)
	require.NoError(t, err)

	// Проверяем, что все остальные клиенты получили сообщение
	for i := 1; i < 3; i++ {
		_, receivedMessage, err := clients[i].ReadMessage()
		require.NoError(t, err)

		var received map[string]interface{}
		err = json.Unmarshal(receivedMessage, &received)
		require.NoError(t, err)
		assert.Equal(t, message["content"], received["content"])
	}
}

func TestWebSocketOldMessagesCleanup(t *testing.T) {
	setupTestDB(t)

	server := setupTestServer(t)
	defer server.Close()

	// Создаем старый комментарий
	oldComment := &models.Comment{
		Content:   "old message",
		TopicId:   testTopicID,
		AuthorId:  testUserID,
		CreatedAt: time.Now().AddDate(0, 0, -15), // 15 дней назад
	}
	_, err := models.AddComment(oldComment)
	require.NoError(t, err)

	// Подключаемся к WebSocket
	ws := connectWebSocket(t, server, testTopicID)
	if ws == nil {
		return
	}
	defer ws.Close()

	// Устанавливаем таймаут для чтения
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Получаем историю сообщений
	_, message, err := ws.ReadMessage()
	require.NoError(t, err)

	var messages []models.CommentWithUsername
	err = json.Unmarshal(message, &messages)
	require.NoError(t, err)

	// Проверяем, что старый комментарий был удален
	for _, msg := range messages {
		assert.NotEqual(t, "old message", msg.Content)
	}
}

func TestWebSocketInvalidMessage(t *testing.T) {
	setupTestDB(t)

	server := setupTestServer(t)
	defer server.Close()

	// Подключаемся к WebSocket
	ws := connectWebSocket(t, server, testTopicID)
	if ws == nil {
		return
	}
	defer ws.Close()

	// Устанавливаем таймаут для чтения
	ws.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Пропускаем историю сообщений
	_, _, err := ws.ReadMessage()
	require.NoError(t, err)

	// Отправляем неверное сообщение
	invalidMessage := []byte("invalid json")
	err = ws.WriteMessage(gorilla.TextMessage, invalidMessage)
	require.NoError(t, err)

	// Проверяем, что соединение все еще активно
	err = ws.WriteMessage(gorilla.PingMessage, nil)
	require.NoError(t, err)
}
