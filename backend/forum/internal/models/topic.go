package models

import (
	"database/sql"
	"fmt"
	"forum/backend/forum/db"
	"log"
	"time"
)

type Topic struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	AuthorId    int            `json:"author_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func GetAllTopics() []Topic {
	var topics []Topic

	rows, err := db.Db.Query("SELECT * FROM topics")
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var t Topic

		err := rows.Scan(&t.ID, &t.Title, &t.Description,
			&t.AuthorId, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			log.Println("Ошибка при сканировании строки:", err)
			continue
		}
		topics = append(topics, t)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка после итерации по строкам:", err)
	}
	return topics
}

func GetTopicByID(id int) (*Topic, error) {
	var t Topic
	err := db.Db.QueryRow("SELECT * FROM topics WHERE id = $1", id).
		Scan(&t.ID, &t.Title, &t.Description,
			&t.AuthorId, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func AddTopic(t *Topic) (int, error) {
	var id int
	query := `
		INSERT INTO topics (title, description, author_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := db.Db.QueryRow(query, t.Title, t.Description, t.AuthorId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить топик: %w", err)
	}
	return id, nil
}

func DeleteTopicByID(id int) error {
	query := `DELETE FROM topics WHERE id = $1`
	_, err := db.Db.Exec(query, id)
	if err != nil {
		log.Println("!!!")
		return fmt.Errorf("не удалось удалить топик: %w", err)
	}
	return nil
}

func PutTopic(id int, updated *Topic) (Topic, error) {
	query := `
		UPDATE topics 
		SET title = $1, description = $2
		WHERE id = $3
		RETURNING id, title, description, author_id, created_at, updated_at
	`
	log.Println(updated.Title)
	log.Println(updated.Description.String)
	var t Topic
	err := db.Db.QueryRow(query, updated.Title, updated.Description.String, id).Scan(
		&t.ID, &t.Title, &t.Description.String, &t.AuthorId, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return Topic{}, fmt.Errorf("не удалось обновить топик: %w", err)
	}
	return t, nil
}
