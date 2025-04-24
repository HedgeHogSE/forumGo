package handlers

import (
	"auth/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllUsers())
}

func GetUser(c *gin.Context) {
	// Проверим валидность ID
	if userID, err := strconv.Atoi(c.Param("user_id")); err == nil {
		// Проверим существование топика
		if user, err := models.GetUserByID(userID); err == nil {
			c.JSON(http.StatusOK, user)

		} else {
			log.Println(err)
			// Если топика нет, прервём с ошибкой
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		log.Println(err)
		// При некорректном ID в URL, прервём с ошибкой
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
