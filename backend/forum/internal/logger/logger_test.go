package logger_test

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/HedgeHogSE/forum/backend/forum/internal/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestInitLogger(t *testing.T) {
	// Сохраняем оригинальный stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Создаем pipe для перехвата вывода
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Инициализируем логгер
	logger.InitLogger()

	// Проверяем, что уровень логирования установлен правильно
	assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())

	// Проверяем формат времени
	assert.Equal(t, time.RFC3339, zerolog.TimeFieldFormat)

	// Логируем тестовое сообщение
	log.Info().Msg("test message")

	// Даем время на запись сообщения
	time.Sleep(100 * time.Millisecond)

	// Закрываем pipe для записи
	w.Close()

	// Читаем вывод
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, что вывод содержит ожидаемые строки
	assert.Contains(t, output, "INF")
	assert.Contains(t, output, "test message")
}

func TestGetLogger(t *testing.T) {
	// Сохраняем оригинальный stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Создаем pipe для перехвата вывода
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Инициализируем логгер
	logger.InitLogger()

	// Получаем логгер для компонента
	componentLogger := logger.GetLogger("test_component")

	// Проверяем, что логгер создан
	assert.NotNil(t, componentLogger)

	// Логируем тестовое сообщение
	componentLogger.Info().Msg("test message")

	// Даем время на запись сообщения
	time.Sleep(100 * time.Millisecond)

	// Закрываем pipe для записи
	w.Close()

	// Читаем вывод
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, что вывод содержит ожидаемые строки
	assert.Contains(t, output, "test_component")
	assert.Contains(t, output, "test message")
	assert.Contains(t, output, "INF")
}

func TestLoggerLevels(t *testing.T) {
	// Сохраняем оригинальный stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Создаем pipe для перехвата вывода
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Инициализируем логгер
	logger.InitLogger()

	// Получаем логгер для компонента
	componentLogger := logger.GetLogger("test_component")

	// Тестируем разные уровни логирования
	componentLogger.Debug().Msg("debug message")
	componentLogger.Info().Msg("info message")
	componentLogger.Warn().Msg("warn message")
	componentLogger.Error().Msg("error message")

	// Даем время на запись сообщений
	time.Sleep(100 * time.Millisecond)

	// Закрываем pipe для записи
	w.Close()

	// Читаем вывод
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, что вывод содержит ожидаемые строки
	assert.NotContains(t, output, "debug message") // Debug не должен быть виден
	assert.Contains(t, output, "info message")
	assert.Contains(t, output, "warn message")
	assert.Contains(t, output, "error message")
}
