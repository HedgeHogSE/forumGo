package handlers

import (
	"database/sql"
	"forum/backend/external"
	"forum/backend/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetAllTopicsWithUsername(c *gin.Context) {
	type TopicWithUser struct {
		ID          int            `json:"id"`
		Title       string         `json:"title"`
		Description sql.NullString `json:"description"`
		Name        string         `json:"name"`
	}

	topics := make([]TopicWithUser, 0)

	topics1 := models.GetAllTopics()

	for i := 0; i < len(topics1); i++ {
		name, err := external.GetUsernameFromAuth(topics1[i].AuthorId)
		if err != nil {
			panic(err)
		}
		topics = append(topics, TopicWithUser{
			ID:          topics1[i].ID,
			Title:       topics1[i].Title,
			Description: topics1[i].Description,
			Name:        name,
		})
	}
	c.JSON(http.StatusOK, topics)
}

func GetTopicWithData(c *gin.Context) {
	type TopicWithData struct {
		ID          int              `json:"id"`
		Title       string           `json:"title"`
		Description sql.NullString   `json:"description"`
		CreatedAt   time.Time        `json:"created_at"`
		Username    string           `json:"username"`
		Comments    []models.Comment `json:"comments"`
		AuthorID    int              `json:"author_id"`
	}

	if topicID, err := strconv.Atoi(c.Param("topic_id")); err == nil {

		if topic, err := models.GetTopicByID(topicID); err == nil {
			username, err := external.GetUsernameFromAuth(topic.AuthorId)
			if err != nil {
				log.Println("ошибка при получении имени:", err)
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "auth service unavailable"})
				return
			}
			comments, err1 := models.GetCommentsByTopicID(topicID)
			if err1 != nil {
				panic(err1)
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

			c.JSON(http.StatusOK, res)

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

func GetTopic(c *gin.Context) {
	// Проверим валидность ID
	if topicID, err := strconv.Atoi(c.Param("topic_id")); err == nil {
		// Проверим существование топика
		if topic, err := models.GetTopicByID(topicID); err == nil {
			c.JSON(http.StatusOK, topic)

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

func PostNewTopic(c *gin.Context) {
	type CreateTopicInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		AuthorId    int    `json:"author_id"`
	}
	var newTopic CreateTopicInput
	if err := c.ShouldBindJSON(&newTopic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.AddTopic(&models.Topic{Title: newTopic.Title,
		Description: sql.NullString{
			String: newTopic.Description,
			Valid:  newTopic.Description != "",
		},
		AuthorId: newTopic.AuthorId})
	c.JSON(http.StatusCreated, newTopic)
}

func DeleteTopic(c *gin.Context) {
	if topicID, err := strconv.Atoi(c.Param("topic_id")); err == nil {
		models.DeleteTopicByID(topicID)
		c.JSON(http.StatusNoContent, gin.H{
			"message": "Delete Article",
		})
		return

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func PutTopic(c *gin.Context) {
	type UpdateTopicInput struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	id, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	var newTopic UpdateTopicInput
	if err := c.ShouldBindJSON(&newTopic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := models.PutTopic(id, &models.Topic{
		Title: newTopic.Title,
		Description: sql.NullString{
			String: newTopic.Description,
			Valid:  newTopic.Description != "",
		},
	})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}
