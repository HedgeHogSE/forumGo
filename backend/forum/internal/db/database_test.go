package db_test

import (
	"os"
	"testing"

	"github.com/HedgeHogSE/forum/backend/forum/internal/db"
)

func TestSetupDB(t *testing.T) {
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

	// Тестируем подключение к БД
	err := db.SetupDB()
	if err != nil {
		t.Errorf("SetupDB вернул ошибку: %v", err)
	}

	// Проверяем, что подключение установлено
	if db.Db == nil {
		t.Error("Db не инициализирован")
	}

	// Проверяем, что можем выполнить простой запрос
	_, err = db.Db.Exec("SELECT 1")
	if err != nil {
		t.Errorf("Не удалось выполнить тестовый запрос: %v", err)
	}

	// Закрываем соединение
	db.Db.Close()
}

func TestSetupDB_InvalidCredentials(t *testing.T) {
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

	// Устанавливаем неверные учетные данные
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "invalid_user")
	os.Setenv("DB_PASSWORD", "invalid_password")
	os.Setenv("DB_NAME", "invalid_db")
	os.Setenv("DB_SSLMODE", "disable")

	// Проверяем, что функция возвращает ошибку
	err := db.SetupDB()
	if err == nil {
		t.Error("SetupDB не вернул ошибку при неверных учетных данных")
	}
}

func TestSetupDB_InvalidHost(t *testing.T) {
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

	// Устанавливаем неверный хост
	os.Setenv("DB_HOST", "invalid_host")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "Sashaezhak2006")
	os.Setenv("DB_NAME", "forum_test")
	os.Setenv("DB_SSLMODE", "disable")

	// Проверяем, что функция возвращает ошибку
	err := db.SetupDB()
	if err == nil {
		t.Error("SetupDB не вернул ошибку при неверном хосте")
	}
}
