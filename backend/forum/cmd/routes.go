package main

import (
	handlers2 "forum/backend/forum/handlers"
	"forum/backend/forum/middleware"
	"log"
)

func initializeRoutes() {
	log.Println("Initializing routes")
	router.Use(middleware.CorsMiddleware())

	topicRoutes := router.Group("/topics")
	{
		topicRoutes.GET("", handlers2.GetAllTopicsWithUsername)
		topicRoutes.GET("/:topic_id", handlers2.GetTopicWithData)
		topicRoutes.POST("", handlers2.PostNewTopic)
		topicRoutes.DELETE("/:topic_id", handlers2.DeleteTopic)
		topicRoutes.PUT("/:topic_id", handlers2.PutTopic)
	}

	commentRoutes := router.Group("/comments")
	{
		commentRoutes.GET("", handlers2.GetAllComments)
		commentRoutes.GET("/:comment_id", handlers2.GetComment)
		commentRoutes.POST("", handlers2.PostNewComment)
		commentRoutes.DELETE("/:comment_id", handlers2.DeleteComment)
		commentRoutes.PUT("/:comment_id", handlers2.PutComment)
	}
}
