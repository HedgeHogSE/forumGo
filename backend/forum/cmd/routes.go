package main

import (
	"forum/backend/forum/internal/handlers"
	"forum/backend/forum/internal/middleware"
	"forum/backend/forum/internal/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func initializeRoutes() {
	log.Println("Initializing routes")
	router.Use(middleware.CorsMiddleware())

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		websocket.HandleConnections(c.Writer, c.Request)
	})

	topicRoutes := router.Group("/topics")
	{
		topicRoutes.GET("", handlers.GetAllTopicsWithUsername)
		topicRoutes.GET("/:topic_id", handlers.GetTopicWithData)
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
