package models

import (
	"database/sql"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateChat создает новый чат и добавляет участников
func CreateChat(db *sql.DB, chat *Chat, participants []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = tx.QueryRow("INSERT INTO chats (name, created_at) VALUES ($1, $2) RETURNING id",
		chat.Name, chat.CreatedAt).Scan(&chat.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, participantID := range participants {
		_, err = tx.Exec("INSERT INTO chat_participants (chat_id, user_id, joined_at) VALUES ($1, $2, $3)",
			chat.ID, participantID, time.Now())
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetChatsByUserID возвращает все чаты пользователя
func GetChatsByUserID(db *sql.DB, userID string) ([]Chat, error) {
	rows, err := db.Query(`
		SELECT c.id, c.name, c.created_at
		FROM chats c
		JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}
