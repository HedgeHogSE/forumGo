package main

import (
	"backend/db"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	db.SetupDB()

	router = gin.Default()

	initializeRoutes()

	router.Run(":8080")
}
