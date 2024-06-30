package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func GetUserIDsByNicknames(db *sql.DB, nicknames []string) ([]int, error) {
	if len(nicknames) == 0 {
		return nil, errors.New("no nicknames provided")
	}

	query := "SELECT id FROM users WHERE username IN (?" + strings.Repeat(",?", len(nicknames)-1) + ")"
	args := make([]interface{}, len(nicknames))
	for i, nickname := range nicknames {
		args[i] = nickname
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

// CreateChat создает новый чат и добавляет участников
func CreateChat(db *sql.DB, chat *Chat, participants []int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec("INSERT INTO chats (name, created_at) VALUES ($1, $2) RETURNING id", chat.Name, chat.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO chat_participants (chat_id, user_id, joined_at) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, userID := range participants {
		if _, err := stmt.Exec(chatID, userID, time.Now()); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// GetChatsByUserID получает все чаты для конкретного пользователя
func GetChatsByUserID(db *sql.DB, userID string) ([]Chat, error) {
	query := `
        SELECT c.id, c.name, c.created_at
        FROM chats c
        JOIN chat_participants cp ON c.id = cp.chat_id
        WHERE cp.user_id = $1`

	rows, err := db.Query(query, userID)
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
