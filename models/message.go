package models

import (
	"backMessage/database"
	"log"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func SendMessage(chatID, userID int, content string) error {
	query := `INSERT INTO messages (chat_id, user_id, content, created_at) VALUES ($1, $2, $3, $4)`
	_, err := database.DB.Exec(query, chatID, userID, content, time.Now())
	if err != nil {
		log.Printf("Error inserting message into database: %v", err)
		return err
	}

	log.Printf("Message sent to chat %d", chatID)
	return nil
}
