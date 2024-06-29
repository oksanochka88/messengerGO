package models

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type UserToken struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

func SaveToken(userID, token string) error {
	connStr := "postgres://postgres:roman@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO user_tokens (user_id, token, created_at) VALUES ($1, $2, $3)",
		userID, token, time.Now())
	return err
}

func DeleteToken(db *sql.DB, token string) error {
	_, err := db.Exec("DELETE FROM user_tokens WHERE token = $1", token)
	return err
}

func IsTokenValid(token string) (bool, error) {
	connStr := "postgres://postgres:roman@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var createdAt time.Time
	err = db.QueryRow("SELECT created_at FROM user_tokens WHERE token = $1", token).Scan(&createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("token not found")
		}
		return false, err
	}

	// Example of token expiration after 24 hours
	if time.Since(createdAt) > 24*time.Hour {
		return false, errors.New("token expired")
	}

	return true, nil
}

func GetTokensByUserID(userID string) ([]UserToken, error) {
	connStr := "postgres://postgres:roman@localhost/postgres?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, user_id, token, created_at FROM user_tokens WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []UserToken
	for rows.Next() {
		var token UserToken
		if err := rows.Scan(&token.ID, &token.UserID, &token.Token, &token.CreatedAt); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}
