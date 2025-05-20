package internal

import (
	"forum/backend/auth/internal/handlers"
	"forum/backend/auth/internal/middleware"
	"log"
)

func initializeRoutes() {
	log.Println("Initializing routes")
	router.Use(middleware.CorsMiddleware())

	// Публичные маршруты
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", handlers.Login)
		authRoutes.POST("/register", handlers.Register)
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(middleware.AuthMiddleware())
	{
		userRoutes.GET("", handlers.GetAllUsers)
		userRoutes.GET("/:user_id", handlers.GetUser)
		userRoutes.POST("", handlers.PostNewUser)
		userRoutes.DELETE("/:user_id", handlers.DeleteUser)
		userRoutes.PUT("/:user_id", handlers.PutUser)
	}
}
