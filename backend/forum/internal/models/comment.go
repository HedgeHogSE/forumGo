package models

import (
	"fmt"
	"log"
	"time"

	"github.com/HedgeHogSE/forum/backend/forum/internal/db"
	"github.com/HedgeHogSE/forum/backend/forum/internal/external"
)

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	AuthorId  int       `json:"author_id"`
	TopicId   int       `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentWithUsername struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	AuthorId  int       `json:"author_id"`
	TopicId   int       `json:"topic_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Username  string    `json:"username"`
}

// CommentService реализует интерфейс grpc.CommentService
type CommentService struct{}

// NewCommentService создает новый экземпляр CommentService
func NewCommentService() *CommentService {
	return &CommentService{}
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

// GetCommentsByAuthorID получает комментарии по ID автора
func GetCommentsByAuthorID(authorID int) ([]Comment, error) {
	var comments []Comment
	rows, err := db.Db.Query("SELECT * FROM comments WHERE author_id = $1", authorID)
	if err != nil {
		return nil, err
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
	return comments, nil
}

// GetCommentsByAuthorID получает комментарии по ID автора
func (s *CommentService) GetCommentsByAuthorID(authorID int) ([]Comment, error) {
	return GetCommentsByAuthorID(authorID)
}

func GetCommentsByTopicID(id int) ([]CommentWithUsername, error) {
	var comments []CommentWithUsername
	rows, err := db.Db.Query("SELECT * FROM comments WHERE topic_id = $1 ORDER BY created_at", id)
	if err != nil {
		log.Println(rows)
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var c CommentWithUsername
		err := rows.Scan(&c.ID, &c.Content, &c.AuthorId, &c.TopicId, &c.CreatedAt, &c.UpdatedAt)
		username, err := external.GetUsernameByUserID(c.AuthorId)
		if err != nil {
			log.Println("Ошибка при сканировании строки:", err)
			continue
		}
		c.Username = username
		comments = append(comments, c)
	}

	if err != nil {
		log.Println("Ошибка при получении имени пользователя:", err)
		return nil, err
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка после итерации по строкам:", err)
	}
	return comments, nil
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
	result, err := db.Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить комментарий: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества удаленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("комментарий с id %d не найден", id)
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
