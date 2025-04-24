package handlers

import (
	"backend/models"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"

	"context"
	"forum/protos/go"
	"github.com/gin-gonic/gin"
)

func GetAllTopics(c *gin.Context) {
	c.JSON(http.StatusOK, models.GetAllTopics())
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
	var newTopic models.Topic
	if err := c.ShouldBindJSON(&newTopic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	models.AddTopic(&models.Topic{Title: newTopic.Title, Description: newTopic.Description, AuthorId: newTopic.AuthorId})
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
	id, err := strconv.Atoi(c.Param("topic_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid article id"})
		return
	}

	var newTopic models.Topic
	if err := c.ShouldBindJSON(&newTopic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := models.PutTopic(id, newTopic)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func getUsernameFromAuth(userID int) (string, error) {
	conn, err := grpc.Dial("auth-service:50051", grpc.WithInsecure())
	if err != nil {
		return "", fmt.Errorf("failed to connect to auth service: %v", err)
	}
	defer conn.Close()

	client := proto.NewAuthServiceClient(conn)
	req := &proto.UserIDRequest{UserId: int32(userID)}
	stream, err := client.GetUsernameByUserID(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to get username from auth service: %v", err)
	}

	// Получаем имя пользователя
	res, err := stream.Recv()
	if err != nil {
		return "", fmt.Errorf("error receiving response: %v", err)
	}

	return res.GetUsername(), nil
}
