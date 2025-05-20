package main

import (
	"forum/backend/forum/db"
	"forum/backend/forum/grpc"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	db.SetupDB()

	router = gin.Default()

	initializeRoutes()

	// Запускаем gRPC сервер в отдельной горутине
	go grpc.StartGRPCServer()

	router.Run(":8080")
}
