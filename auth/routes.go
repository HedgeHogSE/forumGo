package main

import (
	"auth/handlers"
	"auth/middleware"
	"log"
)

func initializeRoutes() {
	log.Println("Initializing routes")
	router.Use(middleware.CorsMiddleware())

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("", handlers.GetAllUsers)
		userRoutes.GET("/:user_id", handlers.GetUser)
		userRoutes.POST("", handlers.PostNewUser)
		userRoutes.DELETE("/:user_id", handlers.DeleteUser)
		userRoutes.PUT("/:user_id", handlers.PutUser)
	}
}
