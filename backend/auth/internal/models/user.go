package models

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/backend/auth/db"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetAllUsers() []User {
	var users []User

	rows, err := db.Db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal("Ошибка при выполнении запроса:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u User

		err := rows.Scan(&u.ID, &u.Name, &u.Username,
			&u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			log.Println("Ошибка при сканировании строки:", err)
			continue
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка после итерации по строкам:", err)
	}
	return users
}

func GetUserByID(id int) (*User, error) {
	var u User
	err := db.Db.QueryRow("SELECT * FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Username,
			&u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func AddUser(u *User) (int, error) {
	var id int
	query := `
		INSERT INTO users (name, username, email, password_hash, is_admin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`

	err := db.Db.QueryRow(query, &u.Name, &u.Username, &u.Email, &u.PasswordHash, &u.IsAdmin).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("не удалось добавить пользователя: %w", err)
	}
	return id, nil
}

func DeleteUserByID(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("не удалось удалить пользователя: %w", err)
	}
	return nil
}

func PutUser(id int, updated User) (User, error) {
	query := `
		UPDATE users 
		SET name = $1, username = $2, email = $3, password_hash = $4, is_admin = $5
		WHERE id = $6
		RETURNING id, name, username, email, password_hash, is_admin, created_at, updated_at
	`
	var u User
	err := db.Db.QueryRow(query, updated.Name, updated.Username, updated.Email,
		updated.PasswordHash, updated.IsAdmin, id).Scan(&u.ID, &u.Name, &u.Username,
		&u.Email, &u.PasswordHash, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return User{}, fmt.Errorf("не удалось обновить пользователя: %w", err)
	}
	return u, nil
}

func GetUsernameByUserID(userID int) (string, error) {
	var username string
	query := "SELECT username FROM users WHERE id = $1"
	err := db.Db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", err
	}
	return username, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthenticateUser(username, password string) (*User, error) {
	var user User
	query := "SELECT * FROM users WHERE username = $1"
	err := db.Db.QueryRow(query, username).Scan(
		&user.ID, &user.Name, &user.Username,
		&user.Email, &user.PasswordHash, &user.IsAdmin,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}

	return &user, nil
}
