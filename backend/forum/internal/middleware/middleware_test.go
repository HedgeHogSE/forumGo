package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HedgeHogSE/forum/backend/forum/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	// Устанавливаем режим тестирования для gin
	gin.SetMode(gin.TestMode)

	// Создаем тестовый роутер
	router := gin.New()
	router.Use(middleware.CorsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Тест 1: Проверка CORS заголовков для обычного запроса
	t.Run("Regular Request", func(t *testing.T) {
		// Создаем тестовый запрос
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)

		// Проверяем статус
		assert.Equal(t, http.StatusOK, w.Code)

		// Проверяем CORS заголовки
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	// Тест 2: Проверка обработки OPTIONS запроса
	t.Run("OPTIONS Request", func(t *testing.T) {
		// Создаем тестовый запрос
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("OPTIONS", "/test", nil)
		router.ServeHTTP(w, req)

		// Проверяем статус
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Проверяем CORS заголовки
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})

	// Тест 3: Проверка с пользовательскими заголовками
	t.Run("Custom Headers", func(t *testing.T) {
		// Создаем тестовый запрос с пользовательскими заголовками
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer token")
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Проверяем статус
		assert.Equal(t, http.StatusOK, w.Code)

		// Проверяем CORS заголовки
		assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
		assert.Equal(t, "GET, POST, PUT, DELETE", w.Header().Get("Access-Control-Allow-Methods"))
		assert.Equal(t, "Content-Type, Authorization", w.Header().Get("Access-Control-Allow-Headers"))
		assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
	})
}
