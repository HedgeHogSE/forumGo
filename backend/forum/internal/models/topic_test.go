package models_test

import (
	"database/sql"
	"testing"

	"github.com/HedgeHogSE/forum/backend/forum/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTopicCRUD(t *testing.T) {
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

	// Тест AddTopic
	t.Run("AddTopic", func(t *testing.T) {
		id, err := models.AddTopic(topic)
		assert.NoError(t, err)
		assert.Greater(t, id, 0)
		topic.ID = id
	})

	// Тест GetTopicByID
	t.Run("GetTopicByID", func(t *testing.T) {
		retrievedTopic, err := models.GetTopicByID(topic.ID)
		assert.NoError(t, err)
		assert.Equal(t, topic.Title, retrievedTopic.Title)
		assert.Equal(t, topic.Description.String, retrievedTopic.Description.String)
		assert.Equal(t, topic.AuthorId, retrievedTopic.AuthorId)
	})

	// Тест GetAllTopics
	t.Run("GetAllTopics", func(t *testing.T) {
		topics := models.GetAllTopics()
		assert.Len(t, topics, 1)
		assert.Equal(t, topic.Title, topics[0].Title)
		assert.Equal(t, topic.Description.String, topics[0].Description.String)
	})

	// Тест PutTopic
	t.Run("PutTopic", func(t *testing.T) {
		updatedTopic := &models.Topic{
			Title:       "Updated Topic",
			Description: sql.NullString{String: "Updated Description", Valid: true},
		}

		result, err := models.PutTopic(topic.ID, updatedTopic)
		assert.NoError(t, err)
		assert.Equal(t, updatedTopic.Title, result.Title)
		assert.Equal(t, updatedTopic.Description.String, result.Description.String)
	})

	// Тест DeleteTopicByID
	t.Run("DeleteTopicByID", func(t *testing.T) {
		err := models.DeleteTopicByID(topic.ID)
		assert.NoError(t, err)

		// Проверяем, что топик удален
		_, err = models.GetTopicByID(topic.ID)
		assert.Error(t, err)
	})
}

func TestGetTopicByID_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся получить несуществующий топик
	_, err := models.GetTopicByID(999)
	assert.Error(t, err)
}

func TestPutTopic_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся обновить несуществующий топик
	updatedTopic := &models.Topic{
		Title:       "Updated Topic",
		Description: sql.NullString{String: "Updated Description", Valid: true},
	}

	_, err := models.PutTopic(999, updatedTopic)
	assert.Error(t, err)
}

func TestDeleteTopicByID_NotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)
	clearTestDB(t)

	// Пытаемся удалить несуществующий топик
	err := models.DeleteTopicByID(999)
	assert.Error(t, err)
}

func TestAddTopic_InvalidData(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Тест с пустым заголовком
	topic := &models.Topic{
		Title:       "",
		Description: sql.NullString{String: "Test Description", Valid: true},
		AuthorId:    testUserID,
	}
	_, err := models.AddTopic(topic)
	if err == nil {
		t.Error("AddTopic не вернул ошибку для пустого заголовка")
	}

	// Тест с несуществующим автором
	topic = &models.Topic{
		Title:       "Test Topic",
		Description: sql.NullString{String: "Test Description", Valid: true},
		AuthorId:    999,
	}
	_, err = models.AddTopic(topic)
	if err == nil {
		t.Error("AddTopic не вернул ошибку для несуществующего автора")
	}

	// Тест с пустым описанием
	topic = &models.Topic{
		Title:       "Test Topic",
		Description: sql.NullString{String: "", Valid: false},
		AuthorId:    testUserID,
	}
	_, err = models.AddTopic(topic)
	if err != nil {
		t.Errorf("AddTopic вернул ошибку для пустого описания: %v", err)
	}
}

func TestPutTopic_InvalidData(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Создаем тестовый топик
	topic := &models.Topic{
		Title:       "Test Topic",
		Description: sql.NullString{String: "Test Description", Valid: true},
		AuthorId:    testUserID,
	}
	id, _ := models.AddTopic(topic)

	// Тест с пустым заголовком
	updated := &models.Topic{
		Title:       "",
		Description: sql.NullString{String: "Updated Description", Valid: true},
	}
	_, err := models.PutTopic(id, updated)
	if err == nil {
		t.Error("PutTopic не вернул ошибку для пустого заголовка")
	}

	// Тест с пустым описанием
	updated = &models.Topic{
		Title:       "Updated Topic",
		Description: sql.NullString{String: "", Valid: false},
	}
	_, err = models.PutTopic(id, updated)
	if err != nil {
		t.Errorf("PutTopic вернул ошибку для пустого описания: %v", err)
	}
}
