package models

import (
	"database/sql"
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
