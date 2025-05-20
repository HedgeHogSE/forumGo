package handlers

import (
	"database/sql"
	"forum/backend/forum/internal/external"
	"forum/backend/forum/internal/logger"
	"forum/backend/forum/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var log1 zerolog.Logger

func init() {
	log1 = logger.GetLogger("topic_handler")
}

func GetAllTopicsWithUsername(c *gin.Context) {
	type TopicWithUser struct {
		ID          int            `json:"id"`
		Title       string         `json:"title"`
		Description sql.NullString `json:"description"`
		Name        string         `json:"name"`
	}

	log1.Info().Msg("Getting all topics with usernames")
	topics := make([]TopicWithUser, 0)
	topics1 := models.GetAllTopics()

	for i := 0; i < len(topics1); i++ {
		name, err := external.GetUsernameFromAuth(topics1[i].AuthorId)
		if err != nil {
			log1.Error().
				Err(err).
				Int("author_id", topics1[i].AuthorId).
				Msg("Failed to get username from auth service")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "auth service unavailable"})
			return
		}
		topics = append(topics, TopicWithUser{
			ID:          topics1[i].ID,
			Title:       topics1[i].Title,
			Description: topics1[i].Description,
			Name:        name,
		})
	}
	log1.Info().Int("topics_count", len(topics)).Msg("Successfully retrieved all topics")
	c.JSON(http.StatusOK, topics)
}

func GetTopicWithData(c *gin.Context) {
	type TopicWithData struct {
		ID          int                          `json:"id"`
		Title       string                       `json:"title"`
		Description sql.NullString               `json:"description"`
		CreatedAt   time.Time                    `json:"created_at"`
		Username    string                       `json:"username"`
		Comments    []models.CommentWithUsername `json:"comments"`
		AuthorID    int                          `json:"author_id"`
	}

	topicID, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		log1.Error().
			Err(err).
			Str("topic_id", c.Param("topic_id")).
			Msg("Invalid topic ID format")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log1.Info().Int("topic_id", topicID).Msg("Getting topic with data")
	topic, err := models.GetTopicByID(topicID)
	if err != nil {
		log1.Error().
			Err(err).
			Int("topic_id", topicID).
			Msg("Failed to get topic")
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	username, err := external.GetUsernameFromAuth(topic.AuthorId)
	if err != nil {
		log1.Error().
			Err(err).
			Int("author_id", topic.AuthorId).
			Msg("Failed to get username from auth service")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
		return
	}

	comments, err := models.GetCommentsByTopicID(topicID)
	if err != nil {
		log1.Error().
			Err(err).
			Int("topic_id", topicID).
			Msg("Failed to get comments for topic")
		panic(err)
	}

	res := TopicWithData{
		ID:          topic.ID,
		Title:       topic.Title,
		Description: topic.Description,
		CreatedAt:   topic.CreatedAt,
		Username:    username,
		Comments:    comments,
		AuthorID:    topic.AuthorId,
	}

	log1.Info().
		Int("topic_id", topicID).
		Int("comments_count", len(comments)).
		Msg("Successfully retrieved topic with data")
	c.JSON(http.StatusOK, res)
}

func GetTopic(c *gin.Context) {
	if topicID, err := strconv.Atoi(c.Param("topic_id")); err == nil {
		if topic, err := models.GetTopicByID(topicID); err == nil {
			c.JSON(http.StatusOK, topic)

		} else {
			log1.Error().
				Err(err).
				Int("topic_id", topicID).
				Msg("Failed to get topic")
			c.AbortWithError(http.StatusNotFound, err)
		}

	} else {
		log1.Error().
			Err(err).
			Str("topic_id", c.Param("topic_id")).
			Msg("Invalid topic ID format")
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PostNewTopic(c *gin.Context) {
	type CreateTopicInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AuthorId    int    `json:"author_id"`
	}

	var newTopic CreateTopicInput
	if err := c.ShouldBindJSON(&newTopic); err != nil {
		log1.Error().
			Err(err).
			Interface("input", newTopic).
			Msg("Invalid topic creation input")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log1.Info().
		Str("title", newTopic.Title).
		Int("author_id", newTopic.AuthorId).
		Msg("Creating new topic")

	topic := &models.Topic{
		Title: newTopic.Title,
		Description: sql.NullString{
			String: newTopic.Description,
			Valid:  newTopic.Description != "",
		},
		AuthorId: newTopic.AuthorId,
	}

	models.AddTopic(topic)
	log1.Info().
		Int("topic_id", topic.ID).
		Msg("Successfully created new topic")
	c.JSON(http.StatusCreated, newTopic)
}

func DeleteTopic(c *gin.Context) {
	topicID, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		log1.Error().
			Err(err).
			Str("topic_id", c.Param("topic_id")).
			Msg("Invalid topic ID format")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	log1.Info().Int("topic_id", topicID).Msg("Deleting topic")
	models.DeleteTopicByID(topicID)
	log1.Info().Int("topic_id", topicID).Msg("Successfully deleted topic")
	c.JSON(http.StatusNoContent, gin.H{
		"message": "Topic deleted",
	})
}

func PutTopic(c *gin.Context) {
	type UpdateTopicInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	id, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		log1.Error().
			Err(err).
			Str("topic_id", c.Param("topic_id")).
			Msg("Invalid topic ID format")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
		return
	}

	var newTopic UpdateTopicInput
	if err := c.ShouldBindJSON(&newTopic); err != nil {
		log1.Error().
			Err(err).
			Interface("input", newTopic).
			Msg("Invalid topic update input")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log1.Info().
		Int("topic_id", id).
		Str("title", newTopic.Title).
		Msg("Updating topic")

	updated, err := models.PutTopic(id, &models.Topic{
		Title: newTopic.Title,
		Description: sql.NullString{
			String: newTopic.Description,
			Valid:  newTopic.Description != "",
		},
	})
	if err != nil {
		log1.Error().
			Err(err).
			Int("topic_id", id).
			Msg("Failed to update topic")
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log1.Info().
		Int("topic_id", id).
		Msg("Successfully updated topic")
	c.JSON(http.StatusOK, updated)
}
