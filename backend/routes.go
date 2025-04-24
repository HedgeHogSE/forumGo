package main

import (
	"backend/handlers"
	"backend/middleware"
	"log"
)

func initializeRoutes() {
	log.Println("Initializing routes")
	router.Use(middleware.CorsMiddleware())

	topicRoutes := router.Group("/topics")
	{
		topicRoutes.GET("", handlers.GetAllTopics)
		topicRoutes.GET("/:topic_id", handlers.GetTopic)
		topicRoutes.POST("", handlers.PostNewTopic)
		topicRoutes.DELETE("/:topic_id", handlers.DeleteTopic)
		topicRoutes.PUT("/:topic_id", handlers.PutTopic)
	}

	commentRoutes := router.Group("/comments")
	{
		commentRoutes.GET("", handlers.GetAllComments)
		commentRoutes.GET("/:comment_id", handlers.GetComment)
		commentRoutes.POST("", handlers.PostNewComment)
		commentRoutes.DELETE("/:comment_id", handlers.DeleteComment)
		commentRoutes.PUT("/:comment_id", handlers.PutComment)
	}
}
