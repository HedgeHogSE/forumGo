package handlers

import (
	"forum/backend/auth/internal/external"
	"forum/backend/auth/internal/logger"
	"forum/backend/auth/internal/models"
	proto "forum/backend/protos/go"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	log = logger.GetLogger("user_handler")
}

func GetAllUsers(c *gin.Context) {
	log.Info().Msg("Getting all users")
	users := models.GetAllUsers()
	log.Info().Int("users_count", len(users)).Msg("Successfully retrieved all users")
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	type UserInfo struct {
		ID        int              `json:"id"`
		Name      string           `json:"name"`
		Username  string           `json:"username"`
		Email     string           `json:"email"`
		IsAdmin   bool             `json:"is_admin"`
		CreatedAt time.Time        `json:"created_at"`
		Comments  []*proto.Comment `json:"comments"`
	}

	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", c.Param("user_id")).
			Msg("Invalid user ID format")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Info().Int("user_id", userID).Msg("Getting user information")
	user, err := models.GetUserByID(userID)
	if err != nil {
		log.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to get user from database")
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	comments, err := external.GetUserCommentsFromBackend(userID)
	if err != nil {
		log.Error().
			Err(err).
			Int("user_id", userID).
			Msg("Failed to get user comments from backend")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user comments"})
		return
	}

	res := UserInfo{
		ID:        userID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt,
		Comments:  comments,
	}

	log.Info().
		Int("user_id", userID).
		Int("comments_count", len(comments)).
		Msg("Successfully retrieved user information")
	c.JSON(http.StatusOK, res)
}

func PostNewUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Error().
			Err(err).
			Interface("user", newUser).
			Msg("Invalid user creation request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Str("username", newUser.Username).
		Str("email", newUser.Email).
		Msg("Creating new user")

	user := &models.User{
		Name:         newUser.Name,
		Username:     newUser.Username,
		Email:        newUser.Email,
		PasswordHash: newUser.PasswordHash,
		IsAdmin:      newUser.IsAdmin,
	}

	models.AddUser(user)
	log.Info().
		Int("user_id", user.ID).
		Str("username", user.Username).
		Msg("Successfully created new user")
	c.JSON(http.StatusCreated, newUser)
}

func DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", c.Param("user_id")).
			Msg("Invalid user ID format")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Info().Int("user_id", userID).Msg("Deleting user")
	models.DeleteUserByID(userID)
	log.Info().Int("user_id", userID).Msg("Successfully deleted user")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "User deleted",
	})
}

func PutUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", c.Param("user_id")).
			Msg("Invalid user ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		log.Error().
			Err(err).
			Interface("user", newUser).
			Msg("Invalid user update request format")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Int("user_id", id).
		Str("username", newUser.Username).
		Msg("Updating user")

	updated, err := models.PutUser(id, newUser)
	if err != nil {
		log.Error().
			Err(err).
			Int("user_id", id).
			Msg("Failed to update user")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Int("user_id", id).
		Msg("Successfully updated user")
	c.JSON(http.StatusOK, updated)
}
