package models

import (
	"database/sql"
	"log"
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateMessage создает новое сообщение в чате
func CreateMessage(db *sql.DB, message *Message) error {
	_, err := db.Exec("INSERT INTO messages (chat_id, user_id, content, created_at) VALUES ($1, $2, $3, $4)",
		message.ChatID, message.UserID, message.Content, message.CreatedAt)
	if err != nil {
		log.Printf("Error creating message: %v", err)
		log.Printf("Message details: ChatID=%d, UserID=%d, Content=%s, CreatedAt=%v", message.ChatID, message.UserID, message.Content, message.CreatedAt)
	}
	return err
}

// GetMessagesByChatID возвращает все сообщения в чате
func GetMessagesByChatID(db *sql.DB, chatID string) ([]Message, error) {
	rows, err := db.Query("SELECT id, chat_id, user_id, content, created_at FROM messages WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.UserID, &message.Content, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func IsUserInChat(db *sql.DB, chatID string, userID string) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1
            FROM chatparticipants
            WHERE chat_id = $1 AND user_id = $2
        )`
	err := db.QueryRow(query, chatID, userID).Scan(&exists)
	return exists, err
}
