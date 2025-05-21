package models_test

import (
	"database/sql"
	"testing"

	"github.com/HedgeHogSE/forum/backend/forum/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCommentCRUD(t *testing.T) {
	// Настраиваем тестовую БД
	setupTestDB(t)
	defer teardownTestDB(t)

	// Очищаем БД перед тестом
	clearTestDB(t)

	// Создаем тестовый топик
	topic := &models.Topic{
		Title:       "Test Topic",
		Description: sql.NullString{String: "Test Description", Valid: true},
		AuthorId:    testUserID,
	}
	topicID, err := models.AddTopic(topic)
	assert.NoError(t, err)

	// Создаем тестовый комментарий
	comment := &models.Comment{
		Content:  "Test Comment",
		AuthorId: testUserID,
		TopicId:  topicID,
	}

	// Тест AddComment
	t.Run("AddComment", func(t *testing.T) {
		id, err := models.AddComment(comment)
		assert.NoError(t, err)
		assert.Greater(t, id, 0)
		comment.ID = id
	})

	// Тест GetCommentByID
	t.Run("GetCommentByID", func(t *testing.T) {
		retrievedComment, err := models.GetCommentByID(comment.ID)
		assert.NoError(t, err)
		assert.Equal(t, comment.Content, retrievedComment.Content)
		assert.Equal(t, comment.AuthorId, retrievedComment.AuthorId)
		assert.Equal(t, comment.TopicId, retrievedComment.TopicId)
	})

	// Тест GetCommentsByTopicID
	t.Run("GetCommentsByTopicID", func(t *testing.T) {
		comments, err := models.GetCommentsByTopicID(topicID)
		assert.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, comment.Content, comments[0].Content)
		assert.Equal(t, comment.AuthorId, comments[0].AuthorId)
		assert.Equal(t, comment.TopicId, comments[0].TopicId)
	})

	// Тест GetCommentsByAuthorID
	t.Run("GetCommentsByAuthorID", func(t *testing.T) {
		comments, err := models.GetCommentsByAuthorID(comment.AuthorId)
		assert.NoError(t, err)
		assert.Len(t, comments, 1)
		assert.Equal(t, comment.Content, comments[0].Content)
		assert.Equal(t, comment.AuthorId, comments[0].AuthorId)
		assert.Equal(t, comment.TopicId, comments[0].TopicId)
	})

	// Тест PutComment
	t.Run("PutComment", func(t *testing.T) {
		updatedComment := models.Comment{
			Content:  "Updated Comment",
			AuthorId: comment.AuthorId,
			TopicId:  comment.TopicId,
		}

		result, err := models.PutComment(comment.ID, updatedComment)
		assert.NoError(t, err)
		assert.Equal(t, updatedComment.Content, result.Content)
		assert.Equal(t, updatedComment.AuthorId, result.AuthorId)
		assert.Equal(t, updatedComment.TopicId, result.TopicId)
	})

	// Тест DeleteCommentByID
	t.Run("DeleteCommentByID", func(t *testing.T) {
		err := models.DeleteCommentByID(comment.ID)
		assert.NoError(t, err)

		// Проверяем, что комментарий удален
		_, err = models.GetCommentByID(comment.ID)
		assert.Error(t, err)
	})
}

func TestGetCommentByID_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся получить несуществующий комментарий
	_, err := models.GetCommentByID(999)
	assert.Error(t, err)
}

func TestGetCommentsByTopicID_Empty(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся получить комментарии для несуществующего топика
	comments, err := models.GetCommentsByTopicID(999)
	assert.NoError(t, err)
	assert.Empty(t, comments)
}

func TestGetCommentsByAuthorID_Empty(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся получить комментарии для несуществующего автора
	comments, err := models.GetCommentsByAuthorID(999)
	assert.NoError(t, err)
	assert.Empty(t, comments)
}

func TestPutComment_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся обновить несуществующий комментарий
	updatedComment := models.Comment{
		Content:  "Updated Comment",
		AuthorId: testUserID,
		TopicId:  1,
	}

	_, err := models.PutComment(999, updatedComment)
	assert.Error(t, err)
}

func TestDeleteCommentByID_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся удалить несуществующий комментарий
	err := models.DeleteCommentByID(999)
	assert.Error(t, err)
}

func TestAddComment_InvalidData(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Тест с пустым содержимым
	comment := &models.Comment{
		Content:  "",
		AuthorId: testUserID,
		TopicId:  1,
	}
	_, err := models.AddComment(comment)
	if err == nil {
		t.Error("AddComment не вернул ошибку для пустого содержимого")
	}

	// Тест с несуществующим автором
	comment = &models.Comment{
		Content:  "Test comment",
		AuthorId: 999,
		TopicId:  1,
	}
	_, err = models.AddComment(comment)
	if err == nil {
		t.Error("AddComment не вернул ошибку для несуществующего автора")
	}

	// Тест с несуществующим топиком
	comment = &models.Comment{
		Content:  "Test comment",
		AuthorId: testUserID,
		TopicId:  999,
	}
	_, err = models.AddComment(comment)
	if err == nil {
		t.Error("AddComment не вернул ошибку для несуществующего топика")
	}
}

func TestPutComment_InvalidData(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Создаем тестовый комментарий
	comment := &models.Comment{
		Content:  "Test comment",
		AuthorId: testUserID,
		TopicId:  1,
	}
	id, _ := models.AddComment(comment)

	// Тест с пустым содержимым
	updated := models.Comment{
		Content:  "",
		AuthorId: testUserID,
		TopicId:  1,
	}
	_, err := models.PutComment(id, updated)
	if err == nil {
		t.Error("PutComment не вернул ошибку для пустого содержимого")
	}

	// Тест с несуществующим автором
	updated = models.Comment{
		Content:  "Updated comment",
		AuthorId: 999,
		TopicId:  1,
	}
	_, err = models.PutComment(id, updated)
	if err == nil {
		t.Error("PutComment не вернул ошибку для несуществующего автора")
	}

	// Тест с несуществующим топиком
	updated = models.Comment{
		Content:  "Updated comment",
		AuthorId: testUserID,
		TopicId:  999,
	}
	_, err = models.PutComment(id, updated)
	if err == nil {
		t.Error("PutComment не вернул ошибку для несуществующего топика")
	}
}
