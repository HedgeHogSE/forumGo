package handlers

import (
	"forum/backend/forum/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllComments(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllComments())
}

func GetComment(c *gin.Context) {
	// Проверим валидность ID
	if commentID, err := strconv.Atoi(c.Param("comment_id")); err == nil {
		// Проверим существование топика
		if comment, err := models.GetCommentByID(commentID); err == nil {
			c.JSON(http.StatusOK, comment)

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

func PostNewComment(c *gin.Context) {
	var newComment models.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.AddComment(&models.Comment{Content: newComment.Content, AuthorId: newComment.AuthorId,
		TopicId: newComment.TopicId})
	c.JSON(http.StatusCreated, newComment)
}

func DeleteComment(c *gin.Context) {
	if commentID, err := strconv.Atoi(c.Param("comment_id")); err == nil {
		models.DeleteCommentByID(commentID)
		c.JSON(http.StatusNoContent, gin.H{
			"message": "Delete Article",
		})
		return

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PutComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	var newComment models.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := models.PutComment(id, newComment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}
