package handlers

import (
	"forum/backend/auth/internal/jwt"
	"forum/backend/auth/internal/logger"
	"forum/backend/auth/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var log1 zerolog.Logger

func init() {
	log1 = logger.GetLogger("auth_handler")
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log1.Error().
			Err(err).
			Interface("request", req).
			Msg("Invalid login request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log1.Info().
		Str("username", req.Username).
		Msg("Attempting user login")

	user, err := models.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		log1.Error().
			Err(err).
			Str("username", req.Username).
			Msg("Authentication failed")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		log1.Error().
			Err(err).
			Int("user_id", user.ID).
			Str("username", user.Username).
			Msg("Failed to generate JWT token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log1.Info().
		Int("user_id", user.ID).
		Str("username", user.Username).
		Msg("User successfully logged in")

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		},
	})
}

func Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log1.Error().
			Err(err).
			Interface("user", newUser).
			Msg("Invalid registration request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log1.Info().
		Str("username", newUser.Username).
		Str("email", newUser.Email).
		Msg("Attempting user registration")

	hashedPassword, err := models.HashPassword(newUser.PasswordHash)
	if err != nil {
		log1.Error().
			Err(err).
			Str("username", newUser.Username).
			Msg("Failed to hash password")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser.PasswordHash = hashedPassword

	userID, err := models.AddUser(&newUser)
	if err != nil {
		log1.Error().
			Err(err).
			Str("username", newUser.Username).
			Msg("Failed to add user to database")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		log1.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to get created user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get created user"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		log1.Error().
			Err(err).
			Int("user_id", user.ID).
			Str("username", user.Username).
			Msg("Failed to generate JWT token for new user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	log1.Info().
		Int("user_id", user.ID).
		Str("username", user.Username).
		Msg("User successfully registered")

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"is_admin": user.IsAdmin,
		},
	})
}
