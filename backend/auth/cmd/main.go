package main

import (
	"forum/backend/auth/internal/db"
	"forum/backend/auth/internal/grpc"
	"forum/backend/auth/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var router *gin.Engine

func main() {
	logger.InitLogger()
	log.Info().Msg("Starting auth service")

	log.Info().Msg("Setting up database connection")
	db.SetupDB()

	log.Info().Msg("Starting gRPC server")
	go grpc.StartGRPCServer()

	log.Info().Msg("Initializing router")
	router = gin.Default()
	InitializeRoutes()

	log.Info().Msg("Starting HTTP server on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start HTTP server")
	}
}
