package handlers

import (
	"forum/backend/auth/external"
	"forum/backend/auth/models"
	proto "forum/backend/protos/go"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllUsers())
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

	// Проверим валидность ID
	if userID, err := strconv.Atoi(c.Param("user_id")); err == nil {
		// Проверим существование пользователя
		if user, err := models.GetUserByID(userID); err == nil {
			comments, err := external.GetUserCommentsFromBackend(userID)
			if err != nil {
				log.Printf("Error getting user comments: %v", err)
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
			c.JSON(http.StatusOK, res)
		} else {
			log.Printf("Error getting user: %v", err)
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		log.Printf("Invalid user ID: %v", err)
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PostNewUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.AddUser(&models.User{Name: newUser.Name, Username: newUser.Username,
		Email: newUser.Email, PasswordHash: newUser.PasswordHash, IsAdmin: newUser.IsAdmin})
	c.JSON(http.StatusCreated, newUser)
}

func DeleteUser(c *gin.Context) {
	if userID, err := strconv.Atoi(c.Param("user_id")); err == nil {
		models.DeleteUserByID(userID)
		c.JSON(http.StatusNoContent, gin.H{
			"message": "Delete User",
		})
		return

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PutUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := models.PutUser(id, newUser)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}
