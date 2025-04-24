package models

import (
	"backend/db"
	"fmt"
	"log"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	AuthorId  int       `json:"author_id"`
	TopicId   int       `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAllComments() []Comment {
	var comments []Comment

	rows, err := db.Db.Query("SELECT * FROM comments")
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment

		err := rows.Scan(&c.ID, &c.Content, &c.AuthorId, &c.TopicId, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			log.Println("Ошибка при сканировании строки:", err)
			continue
		}
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка после итерации по строкам:", err)
	}
	return comments
}

func GetCommentByID(id int) (*Comment, error) {
	var c Comment
	err := db.Db.QueryRow("SELECT * FROM comments WHERE id = $1", id).
		Scan(&c.ID, &c.Content, &c.AuthorId, &c.TopicId, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func AddComment(c *Comment) (int, error) {
	var id int
	query := `
		INSERT INTO comments (content, author_id, topic_id)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err := db.Db.QueryRow(query, c.Content, c.AuthorId, c.TopicId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить комментарий: %w", err)
	}
	return id, nil
}

func DeleteCommentByID(id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := db.Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить комментарий: %w", err)
	}
	return nil
}

func PutComment(id int, updated Comment) (Comment, error) {
	query := `
		UPDATE comments 
		SET content = $1, author_id = $2, topic_id = $3
		WHERE id = $4
		RETURNING id, content, author_id, topic_id, created_at, updated_at
	`
	var c Comment
	err := db.Db.QueryRow(query, updated.Content, updated.AuthorId, updated.TopicId, id).Scan(
		&c.ID, &c.Content, &c.AuthorId, &c.TopicId, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return Comment{}, fmt.Errorf("не удалось обновить комментарий: %w", err)
	}
	return c, nil
}
