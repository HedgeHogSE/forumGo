package main

import (
	"forum/backend/forum/internal/db"
	"forum/backend/forum/internal/grpc"
	"forum/backend/forum/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var router *gin.Engine

func main() {
	logger.InitLogger()
	log.Info().Msg("Starting forum service")

	log.Info().Msg("Setting up database connection")
	db.SetupDB()

	log.Info().Msg("Initializing router")
	router = gin.Default()
	initializeRoutes()

	log.Info().Msg("Starting gRPC server")
	go grpc.StartGRPCServer()

	log.Info().Msg("Starting HTTP server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
