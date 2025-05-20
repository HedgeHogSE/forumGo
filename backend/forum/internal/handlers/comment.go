package handlers

import (
	"forum/backend/forum/internal/logger"
	"forum/backend/forum/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init() {
	log = logger.GetLogger("comment_handler")
}

func GetAllComments(c *gin.Context) {
	log.Info().Msg("Getting all comments")
	comments := models.GetAllComments()
	log.Info().Int("comments_count", len(comments)).Msg("Successfully retrieved all comments")
	c.JSON(http.StatusOK, comments)
}

func GetComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		log.Error().
			Err(err).
			Str("comment_id", c.Param("comment_id")).
			Msg("Invalid comment ID format")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Info().Int("comment_id", commentID).Msg("Getting comment")
	comment, err := models.GetCommentByID(commentID)
	if err != nil {
		log.Error().
			Err(err).
			Int("comment_id", commentID).
			Msg("Failed to get comment")
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	log.Info().Int("comment_id", commentID).Msg("Successfully retrieved comment")
	c.JSON(http.StatusOK, comment)
}

func PostNewComment(c *gin.Context) {
	var newComment models.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		log.Error().
			Err(err).
			Interface("input", newComment).
			Msg("Invalid comment creation input")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Int("topic_id", newComment.TopicId).
		Int("author_id", newComment.AuthorId).
		Msg("Creating new comment")

	comment := &models.Comment{
		Content:  newComment.Content,
		AuthorId: newComment.AuthorId,
		TopicId:  newComment.TopicId,
	}

	models.AddComment(comment)
	log.Info().
		Int("comment_id", comment.ID).
		Int("topic_id", comment.TopicId).
		Msg("Successfully created new comment")
	c.JSON(http.StatusCreated, newComment)
}

func DeleteComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		log.Error().
			Err(err).
			Str("comment_id", c.Param("comment_id")).
			Msg("Invalid comment ID format")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log.Info().Int("comment_id", commentID).Msg("Deleting comment")
	models.DeleteCommentByID(commentID)
	log.Info().Int("comment_id", commentID).Msg("Successfully deleted comment")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Comment deleted",
	})
}

func PutComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("comment_id"))
	if err != nil {
		log.Error().
			Err(err).
			Str("comment_id", c.Param("comment_id")).
			Msg("Invalid comment ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	var newComment models.Comment
	if err := c.ShouldBindJSON(&newComment); err != nil {
		log.Error().
			Err(err).
			Interface("input", newComment).
			Msg("Invalid comment update input")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Int("comment_id", id).
		Int("topic_id", newComment.TopicId).
		Msg("Updating comment")

	updated, err := models.PutComment(id, newComment)
	if err != nil {
		log.Error().
			Err(err).
			Int("comment_id", id).
			Msg("Failed to update comment")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Info().
		Int("comment_id", id).
		Msg("Successfully updated comment")
	c.JSON(http.StatusOK, updated)
}
