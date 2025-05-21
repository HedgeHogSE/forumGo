package models_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/HedgeHogSE/forum/backend/forum/internal/db"
	"github.com/HedgeHogSE/forum/backend/forum/internal/external"
	_ "github.com/lib/pq"
)

var testDB *sql.DB
var testUserID int

// Мок для внешнего сервиса
type mockExternalService struct{}

func (m *mockExternalService) GetUsernameByUserID(userID int) (string, error) {
	if userID == testUserID {
		return "test_user", nil
	}
	return "", fmt.Errorf("user not found")
}

func setupTestDB(t *testing.T) {
	// Сохраняем оригинальные значения переменных окружения
	originalEnv := make(map[string]string)
	for _, key := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
		if value, exists := os.LookupEnv(key); exists {
			originalEnv[key] = value
		}
	}
	defer func() {
		// Восстанавливаем оригинальные значения
		for key, value := range originalEnv {
			os.Setenv(key, value)
		}
	}()

	// Устанавливаем тестовые значения
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "Sashaezhak2006")
	os.Setenv("DB_NAME", "forum_test")
	os.Setenv("DB_SSLMODE", "disable")

	// Подключаемся к тестовой БД
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))

	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Не удалось подключиться к тестовой БД: %v", err)
	}

	// Проверяем подключение
	err = testDB.Ping()
	if err != nil {
		t.Fatalf("Не удалось пинговать тестовую БД: %v", err)
	}

	// Устанавливаем глобальную переменную Db
	db.Db = testDB

	// Создаем тестового пользователя
	createTestUser(t)

	// Устанавливаем мок для внешнего сервиса
	external.GetUsernameByUserID = (&mockExternalService{}).GetUsernameByUserID
}

func createTestUser(t *testing.T) {
	// Проверяем, существует ли уже тестовый пользователь
	var id int
	err := testDB.QueryRow("SELECT id FROM users WHERE username = 'test_user'").Scan(&id)
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
	err = testDB.QueryRow(query).Scan(&testUserID)
	if err != nil {
		t.Fatalf("Не удалось создать тестового пользователя: %v", err)
	}
}

func teardownTestDB(t *testing.T) {
	if testDB != nil {
		testDB.Close()
	}
}

func clearTestDB(t *testing.T) {
	// Очищаем таблицы, но сохраняем тестового пользователя
	_, err := testDB.Exec("TRUNCATE TABLE comments, topics CASCADE")
	if err != nil {
		t.Fatalf("Не удалось очистить тестовую БД: %v", err)
	}
}
